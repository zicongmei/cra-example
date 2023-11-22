package main

import (
	"context"
	"flag"
	"os"
	"path"
	"strings"

	log "github.com/golang/glog"
	rbac "k8s.io/api/rbac/v1"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

const (
	clusterRoleName = "crossplane:provider:provider-anthos-aws-manual:system"
)

var (
	crdDirFlag         = flag.String("crds_dir", "../crds", "directory for CRDs")
	kubeconfigPathFlag = flag.String("kubeconfig", path.Join(os.Getenv("HOME"), ".kube/config"), "path to the kubeconfig")
	serviceAccountFlag = flag.String("service_account", "cloud-resource-accelerator-provider-aws", "service account of the deployment")
	namespaceFlag      = flag.String("namespace", "crossplane-system", "namespace of the deployment")
	cleanFlag          = flag.Bool("clean", false, "clean up the rbac")
)

func main() {
	flag.Set("logtostderr", "true")

	flag.Parse()
	if *crdDirFlag == "" {
		log.Fatal("parameter crds_dir can't be empty")
	}
	if *kubeconfigPathFlag == "" {
		log.Fatal("parameter kubeconfig can't be empty")
	}

	if *cleanFlag {
		deleteRbac(*kubeconfigPathFlag)
	} else {
		crds := getCRDs(*crdDirFlag)
		applyRbac(crds, *kubeconfigPathFlag, *serviceAccountFlag, *namespaceFlag)
	}
}

func getYamlFromDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read dir %q: %v", dir, err)
	}

	var result []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".yaml") {
			result = append(result, e.Name())
		}
	}
	log.Infof("Loaded %d CRD files.", len(result))
	return result
}

func getCRDs(dir string) []*apiext.CustomResourceDefinition {
	files := getYamlFromDir(dir)
	result := []*apiext.CustomResourceDefinition{}
	for i, file := range files {
		file = path.Join(dir, file)
		log.Infof("Handling file #%d %q.", i, file)
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("faile to read file %q: %v", file, err)
		}
		crd := &apiext.CustomResourceDefinition{}

		err = yaml.Unmarshal(content, crd)
		if err != nil {
			log.Fatalf("faile to parse file %q: %v", file, err)
		}
		result = append(result, crd)
	}
	return result
}

func generateRbac(
	crds []*apiext.CustomResourceDefinition,
	serviceAccount string,
	namespace string) (*rbac.ClusterRole, *rbac.ClusterRoleBinding) {
	return generateClusterRole(crds),
		generateClusterRoleBinding(serviceAccount, namespace)
}

func generateClusterRoleBinding(
	serviceAccount string,
	namespace string) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccount,
				Namespace: namespace,
			},
		},
	}
}
func generateClusterRole(crds []*apiext.CustomResourceDefinition) *rbac.ClusterRole {
	clusterRole := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{"", "coordination.k8s.io"},
				Verbs:     []string{"*"},
				Resources: []string{"secrets", "configmaps", "events", "leases"},
			},
		},
	}
	resourceInGroup := map[string][]string{}
	for _, crd := range crds {
		group := crd.Spec.Group
		name := crd.Spec.Names.Plural
		if _, ok := resourceInGroup[group]; !ok {
			resourceInGroup[group] = []string{}
		}
		resourceList := resourceInGroup[group]
		resourceList = append(resourceList, name)
		resourceList = append(resourceList, name+"/status")
		resourceInGroup[group] = resourceList
	}
	for group, resourceList := range resourceInGroup {
		clusterRole.Rules = append(clusterRole.Rules, rbac.PolicyRule{
			APIGroups: []string{group},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"update",
				"patch",
				"create",
			},
			Resources: resourceList,
		})
	}
	return clusterRole
}

func applyRbac(
	crds []*apiext.CustomResourceDefinition,
	kubeconfigPath, serviceAccount, namespace string) {
	clusterRole, clusterRoleBinding := generateRbac(crds, serviceAccount, namespace)

	log.Infof("clusterrole:\n%v", clusterRole)

	client := createKubernetesClient(kubeconfigPath)
	ctx := context.Background()
	if _, err := client.RbacV1().ClusterRoles().
		Create(ctx, clusterRole, metav1.CreateOptions{}); err != nil {
		log.Fatalf("unable to create clusterrole: %v", err)
	}
	if _, err := client.RbacV1().ClusterRoleBindings().
		Create(ctx, clusterRoleBinding, metav1.CreateOptions{}); err != nil {
		log.Fatalf("unable to create clusterrolebing: %v", err)
	}
	log.Infof("Created rbac")
}

func deleteRbac(kubeconfigPath string) {
	client := createKubernetesClient(kubeconfigPath)
	ctx := context.Background()
	if err := client.RbacV1().ClusterRoles().
		Delete(ctx, clusterRoleName, metav1.DeleteOptions{}); err != nil {
		log.Warningf("unable to delete clusterrole: %v", err)
	}
	if err := client.RbacV1().ClusterRoleBindings().
		Delete(ctx, clusterRoleName, metav1.DeleteOptions{}); err != nil {
		log.Warningf("unable to delete clusterrolebing: %v", err)
	}
	log.Infof("Deleted rbac")
}

func createKubernetesClient(kubeconfigPath string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags(
		"", kubeconfigPath,
	)
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("unable to create a client: %v", err)
	}
	return client
}

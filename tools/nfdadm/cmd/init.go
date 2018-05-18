/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type InitCmdFlags struct {
	image      string
	namespace  string
	kubeconfig string
}

var initCmdFlags = &InitCmdFlags{}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initCmdFlags.image, "image", "i", "quay.io/kubernetes_incubator/node-feature-discovery:v0.1.0", "Image to use for the node-feature-discovery binary")
	initCmd.Flags().StringVarP(&initCmdFlags.namespace, "namespace", "n", "default", "Namespace where node-feature-discovery is created")
	initCmd.Flags().StringVarP(&initCmdFlags.kubeconfig, "kubeconfig", "c", defaultKubeconfig(), "Kubeconfig file to use for communicating with the API server")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize NFD",
	Long:  "Initialize Node Feature Discovery on a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Init")

		config, err := clientcmd.BuildConfigFromFlags("", initCmdFlags.kubeconfig)
		if err != nil {
			glog.Exitf("failed to read kubeconfig: %s", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			glog.Exitf("failed to create clientset: %s", err)
		}

		ds, err := createDaemonSet(initCmdFlags, clientset)
		if err != nil {
			glog.Exitf("failed to create DaemonSet: %s", err)
		}

		//fmt.Println(ds.Spec.Template.Spec.Containers[0].Image)
		fmt.Println(ds)
		fmt.Println(err)

	},
}

func createDaemonSet(flags *InitCmdFlags, clientset kubernetes.Interface) (*appsv1.DaemonSet, error) {

	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-feature-discovery",
			Namespace: "default",
			Labels: map[string]string{
				"app": "node-feature-discovery",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "node-feature-discovery",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "node-feature-discovery",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "node-feature-discovery",
							Image: flags.image,
							Args:  []string{"--sleep-interval=60s"},
							Env: []v1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
						},
					},
					HostNetwork: true,
				},
			},
		},
	}

	ds, err := clientset.AppsV1().DaemonSets(flags.namespace).Create(ds)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func defaultKubeconfig() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.Join(usr.HomeDir, ".kube", "config")
}

package k8s

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/buraksezer/olric-kubernetes-plugin/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodAddresses(t *testing.T) {
	cases := []struct {
		Name     string
		Config   *config.Config
		Pods     []corev1.Pod
		Expected []string
	}{
		{
			"Simple pods (no ready, no annotations, etc.)",
			nil,
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase:  corev1.PodRunning,
						PodIP:  "1.2.3.4",
						HostIP: "2.3.4.5",
					},
				},
			},
			[]string{"1.2.3.4"},
		},

		{
			"Simple pods host network",
			&config.Config{HostNetwork: true},
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase:  corev1.PodRunning,
						PodIP:  "1.2.3.4",
						HostIP: "2.3.4.5",
					},
				},
			},
			[]string{"2.3.4.5"},
		},

		{
			"Only running pods",
			nil,
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodPending,
						PodIP: "2.3.4.5",
					},
				},

				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "1.2.3.4",
					},
				},
			},
			[]string{"1.2.3.4"},
		},

		{
			"Only pods that are ready",
			nil,
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodPending,
						PodIP: "2.3.4.5",
					},
				},

				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "ready",
						Conditions: []corev1.PodCondition{
							{
								Type:   corev1.PodReady,
								Status: corev1.ConditionTrue,
							},
						},
					},
				},

				// Not true
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "not-ready",
						Conditions: []corev1.PodCondition{
							{
								Type:   corev1.PodReady,
								Status: corev1.ConditionUnknown,
							},
						},
					},
				},

				// Not ready type, ignored
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "scheduled",
						Conditions: []corev1.PodCondition{
							{
								Type:   corev1.PodScheduled,
								Status: corev1.ConditionUnknown,
							},
						},
					},
				},
			},
			[]string{"ready", "scheduled"},
		},

		{
			"Port annotation (named)",
			nil,
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "1.2.3.4",
					},

					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Ports: []corev1.ContainerPort{
									{
										Name:          "my-port",
										HostPort:      1234,
										ContainerPort: 8500,
									},

									{
										Name:          "http",
										HostPort:      80,
										ContainerPort: 8080,
									},
								},
							},
						},
					},

					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							AnnotationKeyPort: "my-port",
						},
					},
				},
			},
			[]string{"1.2.3.4:8500"},
		},

		{
			"Port annotation (named with host network)",
			&config.Config{HostNetwork: true},
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase:  corev1.PodRunning,
						PodIP:  "1.2.3.4",
						HostIP: "2.3.4.5",
					},

					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										HostPort:      80,
										ContainerPort: 8080,
									},
								},
							},
						},
					},

					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							AnnotationKeyPort: "http",
						},
					},
				},
			},
			[]string{"2.3.4.5:80"},
		},

		{
			"Port annotation (direct)",
			nil,
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
						PodIP: "1.2.3.4",
					},

					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										HostPort:      80,
										ContainerPort: 8080,
									},
								},
							},
						},
					},

					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							AnnotationKeyPort: "4600",
						},
					},
				},
			},
			[]string{"1.2.3.4:4600"},
		},

		{
			"Port annotation (direct with host network)",
			&config.Config{HostNetwork: true},
			[]corev1.Pod{
				{
					Status: corev1.PodStatus{
						Phase:  corev1.PodRunning,
						PodIP:  "1.2.3.4",
						HostIP: "2.3.4.5",
					},

					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										HostPort:      80,
										ContainerPort: 8080,
									},
								},
							},
						},
					},

					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							AnnotationKeyPort: "4600",
						},
					},
				},
			},
			[]string{"2.3.4.5:4600"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			l := log.New(os.Stderr, "", log.LstdFlags)
			addresses, err := PodAddresses(&corev1.PodList{Items: tt.Pods}, tt.Config, l)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if !reflect.DeepEqual(addresses, tt.Expected) {
				t.Fatalf("bad: %#v", addresses)
			}
		})
	}
}

package v1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deployment", func() {
	It("maybe", func() {
		nn := types.NamespacedName{Namespace: "unittest", Name: "deployment"}

		orig := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nn.Name,
				Namespace: nn.Namespace,
				Labels:    map[string]string{"app": "unittest"},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "unittest"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "unittest"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "main",
								Image: "example.com/unittest:latest",
							},
						},
					},
				},
			},
		}
		err := k8sClient.Create(context.TODO(), orig)
		Expect(err).NotTo(HaveOccurred())

		new := &appsv1.Deployment{}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(context.TODO(), nn, new)
			Expect(err).NotTo(HaveOccurred())

			Expect(new.Annotations).ToNot(BeNil())
			value, ok := new.Annotations["example-mutating-admission-webhook"]
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("true"))
		}, time.Second*10, time.Millisecond*250).Should(Succeed())
	})
})

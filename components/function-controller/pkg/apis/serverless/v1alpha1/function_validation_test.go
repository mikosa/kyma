package v1alpha1

import (
	"context"
	"os"
	"testing"

	"github.com/vrischmann/envconfig"

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFunctionSpec_validateResources(t *testing.T) {
	minusOne := int32(-1)
	zero := int32(0)
	one := int32(1)
	os.Setenv("RESERVED_ENVS", "K_CONFIGURATION")
	os.Setenv("MIN_REPLICAS_VALUE", "1")

	for testName, testData := range map[string]struct {
		givenFunc              Function
		expectedError          gomega.OmegaMatcher
		specifiedExpectedError gomega.OmegaMatcher
	}{
		"Should be ok": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source:      "test-source",
					MinReplicas: &one,
					MaxReplicas: &one,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("50m"),
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
					},
				},
			},
			expectedError: gomega.BeNil(),
		},
		"Should validate all fields without error": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source: "test-source",
					Deps:   " { test }     \t\n",
					Env: []corev1.EnvVar{
						{
							Name:  "test",
							Value: "test",
						},
						{
							Name:  "config",
							Value: "test",
						},
					},
					Labels: map[string]string{
						"shoul-be-ok": "test",
						"test":        "test",
					},
					MinReplicas: &one,
					MaxReplicas: &one,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("50m"),
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
					},
				},
			},
			expectedError: gomega.BeNil(),
		},
		"Should return errors on empty function": {
			givenFunc:     Function{},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"metadata.name",
				),
				gomega.ContainSubstring(
					"metadata.namespace",
				),
				gomega.ContainSubstring(
					"spec.source",
				),
			),
		},
		"Should return error on deps validation": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source: "test-source",
					Deps:   "{",
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.deps",
				),
			),
		},
		"Should return error on env validation": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source: "test-source",
					Env: []corev1.EnvVar{
						{
							Name:  "test",
							Value: "test",
						},
						{
							Name:  "K_CONFIGURATION",
							Value: "should reject this",
						},
					},
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.env",
				),
			),
		},
		"Should return error on labels validation": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source: "test-source",
					Labels: map[string]string{
						"shoul-be-ok":      "test",
						"should BE not OK": "test",
					},
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.labels",
				),
			),
		},
		"Should return error on replicas validation": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source:      "test-source",
					MinReplicas: &one,
					MaxReplicas: &minusOne,
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.maxReplicas",
				),
			),
		},
		"Should return error on replicas validation on 0 minReplicas set": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source:      "test-source",
					MinReplicas: &zero, // HPA needs this value to be greater then 0
					MaxReplicas: &one,
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.minReplicas",
				),
			),
		},
		"Should return error on resources validation": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source: "test-source",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("50m"),
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring(
					"spec.resources.limits.cpu",
				),
				gomega.ContainSubstring(
					"spec.resources.limits.memory",
				),
			),
		},
		"should return errors because of minimal config values": {
			givenFunc: Function{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: FunctionSpec{
					Source:      "test-source",
					MinReplicas: &zero,
					MaxReplicas: &zero,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("9m"),
							corev1.ResourceMemory: resource.MustParse("10Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("5m"),
							corev1.ResourceMemory: resource.MustParse("6Mi"),
						},
					},
				},
			},
			expectedError: gomega.HaveOccurred(),
			specifiedExpectedError: gomega.And(
				gomega.ContainSubstring("spec.minReplicas"),
				gomega.ContainSubstring("spec.maxReplicas"),
				gomega.ContainSubstring("spec.resources.requests.cpu"),
				gomega.ContainSubstring("spec.resources.requests.memory"),
				gomega.ContainSubstring("spec.resources.limits.cpu"),
				gomega.ContainSubstring("spec.resources.limits.memory"),
			),
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// given
			g := gomega.NewWithT(t)
			config := &ValidationConfig{}
			envconfig.Init(config)
			ctx := context.WithValue(nil, ValidationConfigKey, *config)

			// when
			errs := testData.givenFunc.Validate(ctx)

			//then
			g.Expect(errs).To(testData.expectedError)
			if testData.specifiedExpectedError != nil {
				g.Expect(errs.Error()).To(testData.specifiedExpectedError)
			}
		})
	}
}

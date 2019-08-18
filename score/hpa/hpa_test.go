package hpa

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestHpaHasTarget(t *testing.T) {
	testcases := []struct {
		hpa           v1.HorizontalPodAutoscaler
		allTargets    []domain.BothMeta
		expectedGrade scorecard.Grade
	}{
		// No match
		{
			hpa: v1.HorizontalPodAutoscaler{
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
			expectedGrade: scorecard.GradeCritical,
		},

		// Match (no namespace)
		{
			hpa: v1.HorizontalPodAutoscaler{
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				},
			},
			expectedGrade: scorecard.GradeAllOK,
		},

		// Match (namespace)
		{
			hpa: v1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace"},
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foospace"},
				},
			},
			expectedGrade: scorecard.GradeAllOK,
		},


		// No match (namespace)
		{
			hpa: v1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace2"},
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foospace"},
				},
			},
			expectedGrade: scorecard.GradeCritical,
		},


		// No match (name)
		{
			hpa: v1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace"},
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "not-foo",
						APIVersion: "apps/v1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foospace"},
				},
			},
			expectedGrade: scorecard.GradeCritical,
		},


		// No match (kind)
		{
			hpa: v1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace"},
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "ReplicaSet",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foospace"},
				},
			},
			expectedGrade: scorecard.GradeCritical,
		},


		// No match (version)
		{
			hpa: v1.HorizontalPodAutoscaler{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace"},
				Spec: v1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: v1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1beta1",
					},
				},
			},
			allTargets: []domain.BothMeta{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1",},
					ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foospace"},
				},
			},
			expectedGrade: scorecard.GradeCritical,
		},
	}

	for _, tc := range testcases {
		fn := hpaHasTarget(tc.allTargets)
		score := fn(tc.hpa)
		assert.Equal(t, tc.expectedGrade, score.Grade)
	}
}

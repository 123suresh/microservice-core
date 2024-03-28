package service

import (
	"context"
	"net/http"

	"github.com/core-api/internal/utils/k8s"
	"github.com/core-api/internal/utils/wordpress"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (svc *Service) CreateWordPress() (string, int, error) {
	wname := "wordpress"
	port := int32(30001)
	err := createWordpressService(wname, port)
	if err == nil {
		wordpress.CreateSecretKey(wname)
		err = wordpress.CreateDatabasePvc(wname)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		err = createWordpressPVC(wname)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		err = wordpress.CreateDatabaseService(wname)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		err = wordpress.CreateDatabaseDeployment(wname)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		err = createWordPressDeployment(wname)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		return "", http.StatusBadRequest, err
	}
	// return err
	return "wordpress created", http.StatusCreated, nil
}

func createWordPressDeployment(wname string) error {
	clientset := k8s.GetConfig()
	namespace := k8s.GetNamespace("wordpress", wname)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: wname,
			Labels: map[string]string{
				"app": wname,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  wname,
					"tier": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  wname,
						"tier": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "wordpress:4.8-apache",
							Name:  wname,
							Env: []corev1.EnvVar{
								{
									Name:  "WORDPRESS_DB_HOST",
									Value: wname + "-mysql",
								},
								{
									Name: "WORDPRESS_DB_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: wname + "-mysql-pass",
											},
											Key: "password",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          wname,
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      wname + "-wordpress-persistent-storage",
									MountPath: "/var/www/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: wname + "-wordpress-persistent-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: wname + "-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}
	log.Info("Creating Wordpress deployment..")
	result, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Wordpress deployment " + result.GetObjectMeta().GetName())
	return nil
}

func int32ptr(i int32) *int32 {
	return &i
}

func createWordpressService(wname string, port int32) error {
	clientset := k8s.GetConfig()
	namespace := k8s.GetNamespace("wordpress", wname)
	nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	// err := CheckIfNamespaceExist(namespace)
	// if err != nil {
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), nsSpec, metav1.CreateOptions{})
	if err != nil {
		log.Error("Failed to create namespace :: ", err)
		return err
	}
	log.Info("Created Namespace " + namespace)
	// }
	servicesClinet := clientset.CoreV1().Services(namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wname,
			Namespace: namespace,
			Labels: map[string]string{
				"app": wname,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     wname,
					Port:     80,
					Protocol: "TCP",
					NodePort: port,
				},
			},
			Selector: map[string]string{
				"app":  wname,
				"tier": "frontend",
			},
			Type: "NodePort",
		},
	}

	log.Info("Creating Wordpress service...")
	_, err = servicesClinet.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Wordpress Service Created Successfully!")
	return nil
}

func createWordpressPVC(pname string) error {
	clinetset := k8s.GetConfig()
	namespace := k8s.GetNamespace("wordpress", pname)

	pvcClinet := clinetset.CoreV1().PersistentVolumeClaims(namespace)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname + "-pv-claim",
			Namespace: namespace,
			Labels: map[string]string{
				"app": pname,
			},
		},
		// Spec: corev1.PersistentVolumeClaimSpec{
		// 	AccessModes: []corev1.PersistentVolumeAccessMode{
		// 		"ReadWriteOnce",
		// 	},
		// 	Resources: corev1.ResourceRequirements{
		// 		Requests: corev1.ResourceList{
		// 			"storage": resource.MustParse("1Gi"),
		// 		},
		// 	},
		// },
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	_, err := pvcClinet.Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Created Wordpress PVC")
	return nil
}

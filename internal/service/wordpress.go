package service

import (
	"context"
	"net/http"

	"github.com/core-api/internal/model"
	"github.com/core-api/internal/utils/k8s"
	"github.com/core-api/internal/utils/wordpress"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (svc *Service) CreateWordPress(req *model.WordPressRequest) (*model.WordPressResponse, int, error) {
	userInput := model.RequestWordPress(req)
	wname := req.Name
	wnamespace := req.Namespace
	count, countErr := svc.repo.CountUser()
	if countErr != nil {
		return nil, http.StatusBadRequest, countErr
	}
	port := int32(30002 + count)
	err := createWordpressService(wname, wnamespace, port)
	if err == nil {
		wordpress.CreateSecretKey(wname, wnamespace)
		err = wordpress.CreateDatabasePvc(wname, wnamespace)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		err = createWordpressPVC(wname, wnamespace)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		err = wordpress.CreateDatabaseService(wname, wnamespace)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		err = wordpress.CreateDatabaseDeployment(wname, wnamespace)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		err = createWordPressDeployment(wname, wnamespace)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		result, err := svc.repo.CreateWordPress(userInput)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		response := result.WordPressResponse()
		return response, http.StatusCreated, nil

	}

	return nil, http.StatusBadRequest, err
}

func createWordPressDeployment(wname string, wnamespace string) error {
	clientset := k8s.GetConfig()
	// namespace := k8s.GetNamespace(wnamespace, wname)
	deploymentsClient := clientset.AppsV1().Deployments(wnamespace)
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

func CheckIfNamespaceExist(namespace string) error {
	client := k8s.GetConfig()
	_, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func createWordpressService(wname string, wnamespace string, port int32) error {
	clientset := k8s.GetConfig()
	// namespace := k8s.GetNamespace(wnamespace, wname)
	nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: wnamespace}}
	err := CheckIfNamespaceExist(wnamespace)
	if err != nil {
		_, err := clientset.CoreV1().Namespaces().Create(context.Background(), nsSpec, metav1.CreateOptions{})
		if err != nil {
			log.Error("Failed to create namespace :: ", err)
			return err
		}
		log.Info("Created Namespace " + wnamespace)
	}
	servicesClinet := clientset.CoreV1().Services(wnamespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wname,
			Namespace: wnamespace,
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

func createWordpressPVC(pname string, wnamespace string) error {
	clinetset := k8s.GetConfig()
	// namespace := k8s.GetNamespace(wnamespace, pname)

	pvcClinet := clinetset.CoreV1().PersistentVolumeClaims(wnamespace)

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname + "-pv-claim",
			Namespace: wnamespace,
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

func (svc *Service) GetWordPress() ([]model.WordPressResponse, int, error) {
	result, err := svc.repo.GetWordPress()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	response := []model.WordPressResponse{}
	for _, detail := range result {
		response = append(response, *detail.WordPressResponse())
	}
	return response, http.StatusOK, nil
}

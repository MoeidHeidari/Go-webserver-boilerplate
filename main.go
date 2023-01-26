package main

import (
	"context"
	"fmt"
	"main/bootstrap"
	"main/docs"
	"main/lib"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// @title SkyFarm
// @description The BEST API you have ever seen
// @host localhost:3000
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	go runSwagger()
	//go run_kubernetes()
	_ = godotenv.Load()
	err := bootstrap.RootApp.Execute()
	if err != nil {
		return
	}

}

func runSwagger() {
	env := lib.NewEnv()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "Skyfarm API"
	r.Run(":" + env.SwaggerPort)
}

func run_kubernetes() {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nodelist, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, n := range nodelist.Items {
		fmt.Println(n)
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		}

		time.Sleep(10 * time.Second)
	}
	// newpod := &corev1.Pod{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "test-pod",
	// 	},
	// 	Spec: corev1.PodSpec{
	// 		Containers: []corev1.Container{
	// 			{Name: "busybox", Image: "busybox:latest", Command: []string{"sleep", "13000000"}},
	// 		},
	// 	},
	// }
	// pod, err := clientset.CoreV1().Pods("default").Create(context.TODO(), newpod, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(pod)
}

package main

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
	var kubeconfig string
	home := homeDir()
	if home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		err := errors.New("kubeconfig not found")
		panic(err)
	}
	fmt.Printf("kubeconfig: %s \n", kubeconfig)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Print(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Print(err)
	}
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()
	// CRD 的做法，是用生成的 client 监听 CRUD 变化。
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Print("AddFunc")
			fmt.Println(reflect.TypeOf(obj))
		}, UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Print("UpdateFunc")
			fmt.Println(reflect.TypeOf(oldObj))
			fmt.Println(reflect.TypeOf(newObj))
		}, DeleteFunc: func(obj interface{}) {
			fmt.Print("DeleteFunc")
			fmt.Println(reflect.TypeOf(obj))
		}})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for{
			time.Sleep(time.Second*5)
			pod, err := podInformer.Lister().Pods("kube-system").Get("csi-cosplugin-external-runner-0")
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println("pod.Status.PodIP：",pod.Status.PodIP)
		}
	}()
	go func(wg *sync.WaitGroup) {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		<-ch
		fmt.Println("time to exit")
		wg.Done()
	}(&wg)
	wg.Wait()
	fmt.Println("end")
}
/*
运行结果：
可以看到，UpdateFunc 产生的推送非常频繁，而 Lister po 则每隔5秒执行一次。程序会在ctrl + C 之后取消无限循环。

UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
pod.Status.PodIP： 172.16.0.16
pod.Status.PodIP： 172.16.0.16
pod.Status.PodIP： 172.16.0.16
pod.Status.PodIP： 172.16.0.16
pod.Status.PodIP： 172.16.0.16
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
pod.Status.PodIP： 172.16.0.16
UpdateFunc*v1.Pod
*v1.Pod
DeleteFunc*v1.Pod
AddFunc*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
UpdateFunc*v1.Pod
*v1.Pod
pod.Status.PodIP： 172.16.0.18
time to exit
end

 */

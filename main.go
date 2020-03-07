package main

import (
    "flag"
    "time"

    "net/http"
    _ "net/http/pprof"


    "github.com/golang/glog"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    // Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
    // _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"


    clientset "github.com/mingregister/crd-learn/pkg/client/clientset/versioned"
    informers "github.com/mingregister/crd-learn/pkg/client/informers/externalversions"
    "github.com/mingregister/crd-learn/pkg/signals"
)

var (
    masterURL  string
    kubeconfig string
)

func main() {
    flag.Parse()

    // 处理信号量
    stopCh := signals.SetupSignalHandler()

    // 处理入参
    cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
    if err != nil {
        glog.Fatalf("Error building kubeconfig: %s", err.Error())
    }

    kubeClient, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
    }

    studentClient, err := clientset.NewForConfig(cfg)
    if err != nil {
        glog.Fatalf("Error building example clientset: %s", err.Error())
    }


    // 调试学习
    go func() {
            http.ListenAndServe(":6060", nil)
    }()

    studentInformerFactory := informers.NewSharedInformerFactory(studentClient, time.Second*30)

    //得到controller
    controller := NewController(kubeClient, studentClient,
        studentInformerFactory.Bolingcavalry().V1().Students())

    //启动informer
    go studentInformerFactory.Start(stopCh)


    //controller开始处理消息
    if err = controller.Run(2, stopCh); err != nil {
        glog.Fatalf("Error running controller: %s", err.Error())
    }
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

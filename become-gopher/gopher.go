package main

import (
	"log"
	"math"
	"syscall/js"
	"sync"
)

var(
	// THREE ...
	THREE js.Value
	once sync.Once
)

// SetupGopher ...
func SetupGopher() {
	THREE = Window.Get("THREE")
	if THREE==js.Undefined() {
		log.Panic("not loaded three.js")
	}
	once.Do(func(){
		world := New()
		world.Render(js.Undefined(), nil)
	})
}

// World ...
type World struct {
	renderer js.Value
	camera js.Value
	scene js.Value
	gopher *Gopher
}

// New ...
func New() *World {
	w := &World{}
	w.renderer = THREE.Get("WebGLRenderer").New(map[string]interface{}{
		"canvas": Document.Call("getElementById", "output"),
		"alpha": true,
	})
	w.scene = THREE.Get("Scene").New()
	w.camera = THREE.Get("PerspectiveCamera").New(45.0, 1.0, 0.1, 1000)
	w.camera.Get("position").Call("set", 0.0, 0.0, 2.0)
	w.scene.Call("add", w.camera)
	w.scene.Call("add", THREE.Get("AmbientLight").New(0xffffff))
	sunlight := THREE.Get("DirectionalLight").New(0xffffff)
	sunlight.Set("radius", 30.0)
	sunlight.Get("position").Call("set", 5, 5, 5)
	w.scene.Call("add", sunlight)
	w.gopher = NewGopher()
	w.scene.Call("add", w.gopher)
	w.OnResize(js.Undefined(), nil)
	Window.Call("addEventListener", "resize", js.FuncOf(w.OnResize))
	Window.Call("addEventListener", "message", js.FuncOf(w.OnMessage))
	return w
}

// OnResize ...
func (w *World) OnResize(this js.Value, args []js.Value) interface{} {
	width := Window.Get("innerWidth").Float();
  height := Window.Get("innerHeight").Float();
  w.renderer.Call("setPixelRatio", Window.Get("devicePixelRatio"));
  w.renderer.Call("setSize", width, height);
  w.camera.Set("aspect", width / height);
	w.camera.Call("updateProjectionMatrix");
	return nil
}

// OnMessage ...
func (w *World) OnMessage(this js.Value, args []js.Value) interface{} {
	ev := args[0]
	poses := []js.Value{}
	data := ev.Get("data")
	for i:=0; i<data.Length(); i++ {
		v := data.Index(i)
		poses = append(poses, v)
	}
	if len(poses)>0 && poses[0].Get("score").Float()>0.4 {
		pose := poses[0].Get("keypoints")
		nose := pose.Index(0).Get("position") // nose
		leye := pose.Index(1).Get("position") // leftEye
		reye := pose.Index(2).Get("position") // rightEye
		sum := nose.Get("y").Float()/600-0.5 + leye.Get("y").Float()/600-0.5 + reye.Get("y").Float()/600-0.5
		w.gopher.Get("position").Call("set", nose.Get("x").Float()/600.0-0.5, -sum/3, 0.0)
	}
	return nil
}

// Render ...
func (w *World) Render(this js.Value, args []js.Value) interface{} {
	w.renderer.Call("render", w.scene, w.camera)
	js.Global().Call("requestAnimationFrame", js.FuncOf(w.Render))
	return nil
}

// Gopher ...
type Gopher struct {
	js.Value
	Head js.Value
	RightArm js.Value
	LeftArm js.Value
}

// NewGopher ...
func NewGopher() *Gopher {
	all := THREE.Get("Group").New()
	head := THREE.Get("Group").New()
	head.Get("position").Call("set", 0, 0.25, 0)
	all.Call("add", head)
	mainMat := THREE.Get("MeshLambertMaterial").New(
		map[string]interface{}{"color": 0x8888ff},
	)
	capsule := js.Global().Call("CapsuleGeometry", 0.5, 1.0, 16)
	capsule.Call("rotateX", math.Pi/2)
	body := THREE.Get("Mesh").New(capsule, mainMat)
	body.Get("position").Call("set", 0, -0.5, 0)
	all.Call("add", body)

	whiteMat := THREE.Get("MeshLambertMaterial").New(
		map[string]interface{}{"color": 0xeeeeee},
	)
	grayMat := THREE.Get("MeshLambertMaterial").New(
		map[string]interface{}{"color": 0x888866},
	)
	blackMat := THREE.Get("MeshLambertMaterial").New(
		map[string]interface{}{"color": 0x444444},
	)
	large := THREE.Get("SphereGeometry").New(0.2, 16, 16)
	small := THREE.Get("SphereGeometry").New(0.05, 8, 8)
	leftEye := THREE.Get("Mesh").New(large, whiteMat)
	iris := THREE.Get("Mesh").New(small, blackMat)
	iris.Get("position").Call("set", 0,0, 0.2, 0.0)
	leftEye.Call("add", iris)
	rightEye := leftEye.Call("clone")
	leftEye.Get("position").Call("set", -0.35, 0.0, 0.3)
	head.Call("add", leftEye)
	rightEye.Get("position").Call("set", +0.35, 0.0, 0.3)
	head.Call("add", rightEye)
	noseGeom := js.Global().Call("CapsuleGeometry", 0.05, 0.05, 16)
	noseGeom.Call("rotateY", math.Pi/2)
	nose := THREE.Get("Mesh").New(noseGeom, blackMat)
	nose.Get("position").Call("set", 0.0, -0.1, 0.5)
	head.Call("add", nose)
	leftLip := THREE.Get("Mesh").New(THREE.Get("SphereGeometry").New(0.1, 8, 8), grayMat)
	rightLip := leftLip.Call("clone")
	leftLip.Get("position").Call("set", -0.08, -0.2, 0.5)
	rightLip.Get("position").Call("set", +0.08, -0.2, 0.5)
	head.Call("add", leftLip)
	head.Call("add", rightLip)
	tooth := THREE.Get("Mesh").New(THREE.Get("BoxGeometry").New(0.1, 0.1, 0.02), whiteMat)
	tooth.Get("position").Call("set", 0.0, -0.3, 0.5)
	head.Call("add", tooth)
	return &Gopher{
		Value: all,
		Head: head,
	}
}


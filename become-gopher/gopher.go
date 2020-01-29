package main

import (
	"log"
	"math"
	"syscall/js"
	"sync"
)

var(
	THREE js.Value
	once sync.Once
)

// SetupGopher ...
func SetupGopher() {
	once.Do(func(){
		var pose js.Value = js.Undefined()
		THREE = Window.Get("THREE")
		log.Println("threejs setup")
		Window.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ev := args[0]
			poses := []js.Value{}
			data := ev.Get("data")
			for i:=0; i<data.Length(); i++ {
				v := data.Index(i)
				poses = append(poses, v)
			}
			if len(poses)>0 && poses[0].Get("score").Float()>0.4 {
				pose = poses[0].Get("keypoints")
			}
			return nil
		}))
		renderer := THREE.Get("WebGLRenderer").New(map[string]interface{}{
			"canvas": Document.Call("getElementById", "output"),
			"alpha": true,
		})
		w, h := 400, 300
		renderer.Call("setSize", w, h)
		scene := THREE.Get("Scene").New()
		js.Global().Set("scene", scene)
		camera := THREE.Get("PerspectiveCamera").New(45.0, w/h, 0.1, 1000)
		//camera.Get("up").Call("set", 0.0, 0.0, 1.0)
		camera.Get("position").Call("set", 0.0, 0.0, 2.0)
		scene.Call("add", camera)
		scene.Call("add", THREE.Get("AmbientLight").New(0xffffff))
		
		sunlight := THREE.Get("DirectionalLight").New(0xffffff)
		sunlight.Set("radius", 30.0)
		sunlight.Get("position").Call("set", 5, 5, 5)
		scene.Call("add", sunlight)
	
		gopher := NewGopher()
		gopher.Get("position").Call("set", 0.0, 0.0, 0.0)
		scene.Call("add", gopher)
		var render js.Func
		render = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if pose!=js.Undefined() {
				nose := pose.Index(0).Get("position") // nose
				leye := pose.Index(1).Get("position") // leftEye
				reye := pose.Index(2).Get("position") // rightEye
				sum := nose.Get("y").Float()/600-0.5 + leye.Get("y").Float()/600-0.5 + reye.Get("y").Float()/600-0.5
				gopher.Get("position").Call("set", nose.Get("x").Float()/600.0-0.5, -sum/3, 0.0)
			}
			renderer.Call("render", scene, camera)
			js.Global().Call("requestAnimationFrame", render)
			return nil
		})
		js.Global().Call("requestAnimationFrame", render)
	})
}
/*
0: {score: 0.9997355341911316, part: "nose", position: {…}}
1: {score: 0.9999076128005981, part: "leftEye", position: {…}}
2: {score: 0.9998937845230103, part: "rightEye", position: {…}}
3: {score: 0.82474285364151, part: "leftEar", position: {…}}
4: {score: 0.968251645565033, part: "rightEar", position: {…}}
5: {score: 0.9511933326721191, part: "leftShoulder", position: {…}}
6: {score: 0.9779354333877563, part: "rightShoulder", position: {…}}
7: {score: 0.02827690728008747, part: "leftElbow", position: {…}}
8: {score: 0.028774429112672806, part: "rightElbow", position: {…}}
9: {score: 0.4714188277721405, part: "leftWrist", position: {…}}
10: {score: 0.47280779480934143, part: "rightWrist", position: {…}}
11: {score: 0.007342240307480097, part: "leftHip", position: {…}}
12: {score: 0.012530015781521797, part: "rightHip", position: {…}}
13: {score: 0.012544909492135048, part: "leftKnee", position: {…}}
14: {score: 0.014684224501252174, part: "rightKnee", position: {…}}
15: {score: 0.007649076171219349, part: "leftAnkle", position: {…}}
16: {score: 0.008718214929103851, part: "rightAnkle", position: {…}}
position{
	x: 0.0,
	y: 0.0
}
*/
type World struct {
	renderer js.Value
	camera js.Value
	scene js.Value
	gopher js.Value
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
	js.Global().Call("requestAnimationFrame", w.Render)
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


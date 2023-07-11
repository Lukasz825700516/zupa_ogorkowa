package main

import (
	// "fmt"
	// "io/ioutil"
	// "net/http"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	// "github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func main() {
	if err := gl.Init(); err != nil {
		println("Failed to initialize opengl ", err)
		os.Exit(1)
	}

	app := adw.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

var vbPane uint32
var ibPane uint32
var shader uint32
var oPlane uint32
func activate(app *adw.Application) {
	bld := gtk.NewBuilderFromFile("./blueprint/main_window.xml")
	win := bld.GetObject("mainWindow").Cast().(*adw.Window)

	surf := bld.GetObject("gl").Cast().(*gtk.GLArea)

	vbPane = 0
	ibPane = 0
	shader = 0
	oPlane = 0

	surf.ConnectResize(func(width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	surf.ConnectRealize(func() {
		var s int32

		v := []float32{
			-1.0, -0.5,
			-1.0, 1.0,
			0.5, -1.0,
			1.0, 1.0,
		}
		i := []int32{
			0, 2, 3,
			0, 3, 1,
		}

		surf.MakeCurrent()

	
		{
			vs := gl.CreateShader(gl.VERTEX_SHADER)
			defer gl.DeleteShader(vs)
			vss, err := ioutil.ReadFile("shaders/vertex_shader.glsl")
			if err != nil {
				println(err)
				os.Exit(1)
			}
			{
				cvss, free := gl.Strs(string(vss))
				defer free()
				gl.ShaderSource(vs, 1, cvss, nil)
				gl.CompileShader(vs)

				gl.GetShaderiv(vs, gl.COMPILE_STATUS, &s)
				if s == 0 {
					var l [512]byte
					gl.GetShaderInfoLog(vs, int32(len(l)), nil, (*uint8)(unsafe.Pointer(&l)))
					println("Vertex shader error ", string(l[:]))
					os.Exit(1)
				}
			}
			fs := gl.CreateShader(gl.FRAGMENT_SHADER)
			defer gl.DeleteShader(fs)
			fss, err := ioutil.ReadFile("shaders/fragment_shader.glsl")
			if err != nil {
				println(err)
				os.Exit(1)
			}
			{
				cfss, free := gl.Strs(string(fss))
				defer free()
				gl.ShaderSource(fs, 1, cfss, nil)
				gl.CompileShader(fs)

				gl.GetShaderiv(fs, gl.COMPILE_STATUS, &s)
				if s == 0 {
					var l [512]byte
					gl.GetShaderInfoLog(fs, int32(len(l)), nil, (*uint8)(unsafe.Pointer(&l)))
					println("Fragment shader error ", string(l[:]))
					os.Exit(1)
				}
			}

			shader = gl.CreateProgram()
			gl.AttachShader(shader, vs)
			gl.AttachShader(shader, fs)
			if err := gl.GetError(); err != gl.NO_ERROR {
				println("gl error attaching ", err)
			}
			gl.LinkProgram(shader)
			if err := gl.GetError(); err != gl.NO_ERROR {
				println("gl error linking ", err)
			}
			println("ok?")
		}

		// gl.GetProgramiv(shader, gl.LINK_STATUS, &s)
			if err := gl.GetError(); err != gl.NO_ERROR {
				println("gl error GetProgramiv ", err)
			}
		if s == 0 {
			var l [512]byte
			gl.GetProgramInfoLog(shader, int32(len(l)), nil, (*uint8)(unsafe.Pointer(&l)))
			println("Linking error1 ", string(l[:]))
			os.Exit(1)
		}


		gl.GenVertexArrays(1, &oPlane)
		gl.GenBuffers(1, &vbPane)
		gl.GenBuffers(1, &ibPane)

		gl.BindVertexArray(oPlane)
		if err := gl.GetError(); err != gl.NO_ERROR {
			println("gl error gen buffers ", err)
		}
		
		gl.BindBuffer(gl.ARRAY_BUFFER, vbPane)
		gl.BufferData(gl.ARRAY_BUFFER, 4 * len(v), gl.Ptr(v), gl.STATIC_DRAW)
		if err := gl.GetError(); err != 0 {
			println("gl error bindd vertex ", err)
		}

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibPane)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4 * len(i), gl.Ptr(i), gl.STATIC_DRAW)
		if err := gl.GetError(); err != gl.NO_ERROR {
			println("gl error bind element ", err)
		}

		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2 * 4, nil)
		gl.EnableVertexAttribArray(0)
		if err := gl.GetError(); err != gl.NO_ERROR {
			println("gl error attrib ", err)
		}
	})
	surf.ConnectRender(func(context gdk.GLContexter) (ok bool) {
		gl.ClearColor(0.1, 0.7, 0.3, 0.7)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shader)
		gl.BindVertexArray(oPlane)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		return true
	})
	surf.ConnectUnrealize(func() {
		gl.DeleteProgram(shader)
		gl.DeleteVertexArrays(1, &oPlane)
		gl.DeleteBuffers(1, &vbPane)
		gl.DeleteBuffers(1, &ibPane)
	})

	app.AddWindow(&win.Window)
	win.Show()
}

package utils

import (
	_ "embed"
	"github.com/thedevsaddam/gojsonq"
	ort "github.com/yalue/onnxruntime_go"
	"image"
	"log"
	"os"
	"strconv"
)
import "fmt"

// 验证码识别
//
//go:embed third_party/onnxruntime.dll
var onnxruntimeDLL []byte

//go:embed third_party/common_old1.onnx
var common_old1 []byte

func writeDLLToDisk() {
	PathExistForCreate("./assets")
	dllPath := "./assets/onnxruntime.dll"
	onnx := "./assets/common_old1.onnx"
	f1 := os.WriteFile(dllPath, onnxruntimeDLL, 0644)
	if f1 != nil {
		log.Fatal(f1)
	}
	f2 := os.WriteFile(onnx, common_old1, 0644)
	if f2 != nil {
		log.Fatal(f2)
	}
}

func AutoVerification(img image.Image) string {
	writeDLLToDisk() // 确保文件都加载了

	ort.SetSharedLibraryPath("./assets/onnxruntime.dll")

	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
	defer ort.DestroyEnvironment()

	img1 := ResizeImage(img, uint(64*img.Bounds().Dx()/img.Bounds().Dy()), 64)
	imgGray := ConvertToGray(img1)

	inputData := ImageToGrayFloatArray(imgGray)
	inputShape := ort.NewShape(1, 1, 64, int64(imgGray.Bounds().Dx()))
	inputTensor, err := ort.NewTensor[float32](inputShape, inputData)

	if err != nil {
		panic(err)
	}

	defer inputTensor.Destroy()
	// This hypothetical network maps a 2x5 input -> 2x3x4 output.
	outputShape := ort.NewShape(1, 18)
	outputTensor, err := ort.NewEmptyTensor[int64](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession("./assets/common_old1.onnx",
		[]string{"input1"}, []string{"output"},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	defer session.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	err = session.Run()
	if err != nil {
		fmt.Errorf(err.Error())
	}

	outputData := outputTensor.GetData()
	codeResult := ""
	for i := 0; i < len(outputData); i++ {
		if outputData[i] != 0 {
			codeResult += gojsonq.New().JSONString(getCharCode()).Find("[" + strconv.Itoa(int(outputData[i])) + "]").(string)
		}
	}
	return codeResult
}

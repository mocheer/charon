package dmap

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/mocheer/charon/res"
	"github.com/mocheer/pluto/calc/kriging"
)

// KrigingTrainArgs
type KrigingTrainArgs struct {
	Values   []float64
	X        []float64
	Y        []float64
	Model    string // model 模型类型
	Bbox     [][2]float64
	CellSize float64
}

// krigingTrainHandle
func krigingGridHandle(c *fiber.Ctx) error {
	args := &KrigingTrainArgs{}
	// req.MustParseArgs(c, args)
	json.Unmarshal([]byte(c.Query("values")), &args.Values)
	json.Unmarshal([]byte(c.Query("x")), &args.X)
	json.Unmarshal([]byte(c.Query("y")), &args.Y)
	json.Unmarshal([]byte(c.Query("bbox")), &args.Bbox)
	//
	args.CellSize = 0.01
	args.Model = "exponential"

	k := kriging.New(args.Values, args.X, args.Y)
	k.Train(args.Model)
	data := k.Grid(args.Bbox, args.CellSize)
	return res.JSON(c, data)
}

// krigingImageHandle
func krigingImageHandle(c *fiber.Ctx) error {
	return nil
}

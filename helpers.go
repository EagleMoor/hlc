package main

import (
	"fmt"
	"strconv"

	"hlc/models"
	"hlc/ptr"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func getParamID(ctx *routing.Context) (int, error) {
	return strconv.Atoi(ctx.Param("id"))
}

func getIntParam(args *fasthttp.Args, name string) (*int, error) {
	i, err := args.GetUint(name)
	if err != nil {
		if err == fasthttp.ErrNoArgValue {
			return nil, nil
		}

		return nil, err
	}
	return ptr.Int(i), nil
}

func getStringParam(args *fasthttp.Args, name string) *string {
	s := args.Peek(name)
	if s == nil {
		return nil
	}

	return ptr.String(string(s))
}

func getGenderParam(args *fasthttp.Args) (*models.UserGender, error) {
	gender := getStringParam(args, "gender")
	if gender == nil {
		return nil, nil
	}

	var g models.UserGender

	if "m" == *gender {
		g = models.UserGenderMale
		return &g, nil
	}
	if "f" == *gender {
		g = models.UserGenderFemale
		return &g, nil
	}

	return nil, fmt.Errorf("wrong gender")
}

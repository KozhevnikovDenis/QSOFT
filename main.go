package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("X-PING", "ping")
		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(Logger())
	router.Delims("{[{", "}]}")
	router.SetFuncMap(template.FuncMap {
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLFiles("./testdata/raw.tmpl")
	
	router.GET("/when/:year", func(c *gin.Context) {
		if _, err := strconv.Atoi(c.Param("year")); err == nil {
			
			t0 := time.Now()
			str_t1, _ := strconv.Atoi(c.Param("year"))
			t1 := time.Date(str_t1, 01, 01, 0, 0, 0, 0, time.UTC)
			days := int(t1.Sub(t0).Hours() / 24)
				
			c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{} {
				"now": time.Now(),
				"day": days,
			})
			
			if (c.MustGet("X-PING").(string) == "ping") {
				c.Set("X-PONG", "pong")
				c.Next()
				log.Println(c.MustGet("X-PONG").(string))
			} else {
				log.Println("X-PING != ping")
			}
			
		} else {
			log.Println(err)
		}
		
	})
	
	router.Run(":8080")
}


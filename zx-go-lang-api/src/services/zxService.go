package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	utils "github.com/GouthamGuna/zx/src/utils"
	"github.com/gin-gonic/gin"
)

/*func StreamJSON(c *gin.Context) {
    filePath := utils.GetRootPath() + "xyz.json" // Path to your large JSON file

    file, err := os.Open(filePath)
    if err != nil {
        c.JSON(404, gin.H{"error": "File not found"})
        log.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)

    // Get pagination parameters from query
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))

    // Get custom query parameters
    customQuery := c.DefaultQuery("query", "")

    // Calculate start and end index
    start := (page - 1) * limit
    end := page * limit

    var totalDocuments int
    var documents []interface{}
    var filteredDocuments []interface{}

    // Read and filter documents based on custom query
    for decoder.More() {
        var data map[string]interface{} // Use a map to handle arbitrary JSON data
        err := decoder.Decode(&data)
        if err != nil {
            log.Println("Error decoding JSON:", err)
            c.JSON(500, gin.H{"error": "Error decoding JSON"})
            return
        }

        // Custom query logic (example: filter documents containing the custom query in the "name" field)
        if customQuery == "" || (customQuery != "" && data["name"] == customQuery) {
            filteredDocuments = append(filteredDocuments, data)
        }
        documents = append(documents, data)
    }

    totalDocuments = len(filteredDocuments)
    totalPages := (totalDocuments + limit - 1) / limit
    remainingPages := totalPages - page

    if start >= totalDocuments {
        c.JSON(200, gin.H{
            "page":            page,
            "limit":           limit,
            "totalDocuments":  totalDocuments,
            "totalPages":      totalPages,
            "remainingPages":  remainingPages,
            "data":            []interface{}{},
        })
        return
    }

    if end > totalDocuments {
        end = totalDocuments
    }

    c.Header("Content-Type", "application/json")
    c.Writer.Write([]byte("[")) // start of json array.

    first := true
    for i := start; i < end; i++ {
        if !first {
            c.Writer.Write([]byte(","))
        } else {
            first = false
        }
        jsonData, err := json.Marshal(filteredDocuments[i])
        if err != nil {
            log.Println("Error marshalling data:", err)
            c.JSON(500, gin.H{"error": "Error marshalling data"})
            return
        }
        c.Writer.Write(jsonData)
    }
    c.Writer.Write([]byte("]")) // end of json array.

    c.JSON(200, gin.H{
        "page":            page,
        "limit":           limit,
        "totalDocuments":  totalDocuments,
        "totalPages":      totalPages,
        "remainingPages":  remainingPages,
    })
}*/

func streamJSON(filename string, c *gin.Context) {
	filePath := utils.GetRootPath() + filename // Path to your large JSON file

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	c.Header("Content-Type", "application/json")
	c.Writer.Write([]byte("[")) // start of json array

	first := true
	for decoder.More() {
		if !first {
			c.Writer.Write([]byte(","))
		} else {
			first = false
		}

		var data interface{} // Or a specific struct
		if err := decoder.Decode(&data); err != nil {
			log.Println("Error decoding JSON:", err)
			c.JSON(500, gin.H{"error": "Error decoding JSON"})
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data:", err)
			c.JSON(500, gin.H{"error": "Error marshalling data"})
			return
		}

		if _, err := c.Writer.Write(jsonData); err != nil {
			log.Println("Error writing JSON to response:", err)
			c.JSON(500, gin.H{"error": "Error writing JSON to response"})
			return
		}
	}

	c.Writer.Write([]byte("]")) // end of json array
}

// GetDataHandler serves the content of a specified JSON file.
func GetDataHandler(c *gin.Context) {
	filename := c.Param("filename")
	streamJSON(filename, c)
}

// GetFilenamesHandler returns a list of filenames with .json suffix
func GetFilenamesHandler(c *gin.Context) {
	filenames := []string{
		"cV9hc2lhbg.json",
		"cV9iaWctZGljaw.json",
		"cV9ob21lbWFkZQ.json",
		"cV9oYXJkY29yZQ.json",
		"cV9qYXBhbmVzZQ.json",
		"cV9rb3JlYW4.json",
		"cV9sYXRpbmE.json",
		"cV9sZXNiaWFucw.json",
		"cV9tYXR1cmU.json",
        "cV9yZWRoZWFk.json",
		"cV9vcmd5.json",
        "cV90ZWVu.json",
		"enguaW1hZ2Vz.json",
	}

	c.JSON(http.StatusOK, gin.H{"size": len(filenames), "filenames": filenames})
}

// GetDataHandler serves the content of a specified JSON file.
// func GetDataHandler(c *gin.Context) {
// 	filename := c.Param("filename")
// 	data, err := LoadData(filename)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or unable to read"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, data)
// }

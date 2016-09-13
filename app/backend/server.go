/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type comment struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var dataFile string

var commentMutex = new(sync.Mutex)

// Handle comments
func handleComments(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	commentMutex.Lock()
	defer commentMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			basepath := path.Dir(dataFile)
			if os.MkdirAll(basepath, 0777) != nil {
				errText := fmt.Sprintf("Unable to create path for the data file (%s)", basepath)
				fmt.Println(errText)
				http.Error(w, errText, http.StatusInternalServerError)
				return
			}
			err := ioutil.WriteFile(dataFile, []byte("[]"), 0644)
			if err != nil {
				errText := fmt.Sprintf("Unable to write comments to data file (%s): %s", dataFile, err)
				fmt.Println(errText)
				http.Error(w, errText, http.StatusInternalServerError)
				return
			} else {
				return
			}

		}
		errText := fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err)
		fmt.Println(errText)
		http.Error(w, errText, http.StatusInternalServerError)
		return
	}

	// Read the comments from the file.
	commentData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		errText := fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err)
		fmt.Println(errText)
		http.Error(w, errText, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var comments []comment
		if err := json.Unmarshal(commentData, &comments); err != nil {
			errorText := fmt.Sprintf("Unable to Unmarshal comments from data file (%s): %s", dataFile, err)
			fmt.Println(errorText)
			http.Error(w, errorText, http.StatusInternalServerError)
			return
		}

		// Add a new comment to the in memory slice of comments
		comments = append(comments, comment{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})

		// Marshal the comments to indented json.
		commentData, err = json.MarshalIndent(comments, "", "    ")
		if err != nil {
			errorText := fmt.Sprintf("Unable to marshal comments to json: %s", err)
			fmt.Println(errorText)
			http.Error(w, errorText, http.StatusInternalServerError)
			return
		}

		// Write out the comments to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, commentData, fi.Mode())
		if err != nil {

			errorText := fmt.Sprintf("Unable to write comments to data file (%s): %s", dataFile, err)
			fmt.Println(errorText)
			http.Error(w, errorText, http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, bytes.NewReader(commentData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(commentData))

	default:
		// Don't know the method, so error
		errText := fmt.Sprintf("Unsupported method: %s", r.Method)
		fmt.Println(errText)
		http.Error(w, errText, http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	dataFile = os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "./data/comments.json"
	}

	http.HandleFunc("/api/comments", handleComments)
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server started: http://" + host + ":" + port + ", dataFile: " + dataFile)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

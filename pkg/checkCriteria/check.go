package gogemcheckcriteria

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elliotchance/orderedmap"
)

/*

Tries to evaluate if all mandatory Links are reachable.
Links are defined relative to the Team main page.

*/
func createAwardMap() *orderedmap.OrderedMap {
	m := orderedmap.NewOrderedMap()

	m.Set("Medals", "#")
	m.Set("Bronze #2 (Attributions)", "Attributions")
	m.Set("Bronze #3 (Project Description)", "Description")
	m.Set("Bronze #4 (Contribution)", "Contribution")
	m.Set("Silver #1 (Engineering Success)", "Engineering")
	m.Set("Silver #2 (Collaboration)", "Collaborations")
	m.Set("Silver #3 (Human Practices)", "Human_Practices")
	m.Set("Silver #4 (Proposed Implementation)", "Implementation")
	m.Set("Gold #1 (Integrated Human Practices)", "Human_Practices")
	m.Set("Gold #3 (Project Modeling)", "Model")
	m.Set("Gold #4 (Proof of Concept)", "Proof_Of_Concept")
	m.Set("Gold #5 (Partnership)", "Partnership")
	m.Set("Gold #6 (Education & Communication)", "Communication")
	m.Set("Special", "#")
	m.Set("Best Education", "Education")
	m.Set("Best Hardware", "Hardware")
	m.Set("Inclusivity Award", "Inclusivity")
	m.Set("Best HP", "Human_Practices")
	m.Set("Best Measurement", "Measurement")
	m.Set("Best Model", "Model")
	m.Set("Best Plant SynBio", "Plant")
	m.Set("Best Software Tool", "Software")
	m.Set("Best Supporting Entrepreneurship", "Entrepreneurship")
	m.Set("Best Sustainable Development Impact", "Sustainable")
	m.Set("Safety and Security Award", "Safety")

	return m
}

func CheckCriteria(team string, year int, url bool) (string, error) {
	result := ""
	var err error
	baseURL := "https://" + fmt.Sprint(year) + ".igem.org/Team:" + team + "/"

	awardMap := createAwardMap()

	for el := awardMap.Front(); el != nil; el = el.Next() {
		medal := el.Key.(string)

		if el.Value.(string) == "#" {
			result += "########################################################\n"
			result += "############################" + medal + "#####################\n"
			result += "########################################################\n"
			continue
		}

		link := baseURL + el.Value.(string)
		resp, err := http.Get(link)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 404 {
			if url {
				result += medal + ": Page Not Found, check your URLs! " + link + "\n"
			} else {
				result += medal + ": Page Not Found, check your URLs!\n"
			}
		} else if resp.StatusCode == 200 {
			byteBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			body := string(byteBody)

			if strings.Contains(body, `cnoarticletext`) || strings.Contains(body, `purged-page-empty`) || strings.Contains(body, `(page does not exist)`){
				if url {
					result += medal + ": Page is Empty, check your URLs! " + link + "\n"
				} else {
					result += medal + ": Page is Empty, check your URLs!\n"
				}
			} else if strings.Contains(body, `judges-will-not-evaluate`) {
				if url {
					result += medal + ": Page is NOT visible to judges, remove the ALERT message! " + link + "\n"
				} else {
					result += medal + ": Page is NOT visible to judges, remove the ALERT message!\n"
				}
			} else {
				if url {
					result += medal + ": Page seems to be okay, check anyways! " + link + "\n"
				} else {
					result += medal + ": Page seems to be okay, check anyways!\n"
				}
			}
		} else {
			if url {
				result += medal + ": Unknown Error! " + fmt.Sprint(resp.StatusCode) + ": " + link + "\n"
			}
			result += medal + ": Unknown Error! " + fmt.Sprint(resp.StatusCode) + "\n"
		}
	}

	return result, err

}

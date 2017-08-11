package jenkins_ci_parser

import (
	"errors"
	"fmt"

	"strings"

	"gopkg.in/yaml.v2"
)

type JenkinsCiYamlParser struct {
	data           []byte
	config         DataBag
	pipelineConfig DataBag
}

func (c *JenkinsCiYamlParser) parseContent() (err error) {
	config := make(DataBag)
	err = yaml.Unmarshal(c.data, config)
	if err != nil {
		return err
	}

	err = config.Sanitize()
	if err != nil {
		return err
	}

	c.config = config

	return
}

func (c *JenkinsCiYamlParser) loadOverall(jenkinsCI *JenkinsCI) (err error) {

	if environments, ok := c.config.GetSlice("environments"); ok {

		jenkinsCI.Environments = make([]Environment, 0)

		s, err := c.getKeyValuePairs(environments)
		if err != nil {
			return errors.New("Must have one environment at least ")
		}

		for _, env := range s {
			jenkinsCI.Environments = append(jenkinsCI.Environments, Environment{
				Name: env.Key,
				Host: env.Value,
			})
		}
	}

	return
}

func (c *JenkinsCiYamlParser) loadPipeline() (err error) {
	pipelineConfig, ok := c.config.GetSubOptions("pipeline")
	if !ok {
		return fmt.Errorf("no pipeline syntax defined")
	}
	c.pipelineConfig = pipelineConfig
	return
}

func (c *JenkinsCiYamlParser) preparePipelineAgent(pipeline *Pipeline) (err error) {
	//if agent, ok := c.pipelineConfig.GetString("agent"); ok {
	//	pipeline.Agent.Mode = c.getTexts()
	//}
	return
}

func (c *JenkinsCiYamlParser) preparePipelineEnv(pipeline *Pipeline) (err error) {
	if environment, ok := c.pipelineConfig.GetString("environment"); ok {

		var scripts Scripts
		scripts, err = c.getCommands(environment)
		if err != nil {
			return
		}
		pipeline.Env = JenkinsEnv{
			Scripts: scripts,
		}
	}

	return
}

func (c *JenkinsCiYamlParser) preparePipelineStages(pipeline *Pipeline) (err error) {

	if stages, ok := c.pipelineConfig.GetSlice("stages"); ok {

		pipeline.Stages = make([]Stage, 0)
		for _, i := range stages {
			s, err := c.getStage(i)
			if err != nil {
				return err
			}

			pipeline.Stages = append(pipeline.Stages, s)

		}
	}

	return
}

func (c *JenkinsCiYamlParser) getStage(commands interface{}) (Stage, error) {
	if command, ok := commands.(map[interface{}]interface{}); ok {
		var name string
		if val, ok := command["name"]; ok {
			if v, ok := val.(string); ok {
				name = v
			}
		}
		var steps Scripts
		var err error
		if val, ok := command["steps"]; ok {
			steps, err = c.getCommands(val)
			if err != nil {
				return Stage{}, err
			}
		}
		if name != "" && len(steps) > 0 {
			stage := Stage{
				Name:  name,
				Steps: steps,
			}
			if val, ok := command["when"]; ok {

				pairs, err := c.getKeyValuePairs(val)

				if err != nil {
					return Stage{}, err
				}
				var whens Whens
				for _, w := range pairs {
					whens = append(whens, When{
						Condition: w.Key,
						Value:     w.Value,
					})
				}
				stage.When = whens
			}

			return stage, nil

		} else {
			return Stage{}, errors.New("Stage must have name and one step at least ")
		}
	}

	return Stage{}, nil
}

func (c *JenkinsCiYamlParser) getCommands(commands interface{}) (Scripts, error) {
	if lines, ok := commands.([]interface{}); ok {
		var steps Scripts
		for _, line := range lines {
			if lineText, ok := line.(string); ok {
				steps = append(steps, lineText)
			} else {
				return Scripts{}, errors.New("unsupported script")
			}
		}
		return steps, nil
	} else if text, ok := commands.(string); ok {
		return Scripts(strings.Split(text, "\n")), nil
	} else if commands != nil {
		return Scripts{}, errors.New("unsupported script")
	}

	return Scripts{}, nil
}

func (c *JenkinsCiYamlParser) getTexts(commands interface{}) ([]string, error) {
	if lines, ok := commands.([]interface{}); ok {
		var text []string
		for _, line := range lines {
			if lineText, ok := line.(string); ok {
				text = append(text, lineText)
			} else {
				return Scripts{}, errors.New("unsupported script")
			}
		}
		return text, nil
	} else if text, ok := commands.(string); ok {
		return Scripts(strings.Split(text, "\n")), nil
	} else if commands != nil {
		return Scripts{}, errors.New("unsupported script")
	}

	return Scripts{}, nil
}

func (c *JenkinsCiYamlParser) getWhens(commands interface{}) (Whens, error) {

	if lines, ok := commands.([]interface{}); ok {
		var whens Whens
		for _, line := range lines {
			if lineText, ok := line.(map[interface{}]interface{}); ok {
				for key, v := range lineText {
					w := When{
						Condition: fmt.Sprintf("%v", key),
						Value:     fmt.Sprintf("%v", v),
					}
					whens = append(whens, w)
				}
			} else {
				return Whens{}, errors.New("unsupported script")
			}
		}
		return whens, nil
	} else if text, ok := commands.(string); ok {
		fmt.Printf("text: %s", text)

		// return Whens(strings.Split(text, ":")), nil
	} else if commands != nil {
		return Whens{}, errors.New("unsupported script")
	}

	return Whens{}, nil
}

func (c *JenkinsCiYamlParser) getKeyValuePairs(commands interface{}) ([]keyValuePair, error) {

	if lines, ok := commands.([]interface{}); ok {
		var pairs []keyValuePair
		for _, line := range lines {
			if lineText, ok := line.(map[interface{}]interface{}); ok {
				for key, v := range lineText {
					w := keyValuePair{
						Key:   fmt.Sprintf("%v", key),
						Value: fmt.Sprintf("%v", v),
					}
					pairs = append(pairs, w)
				}
			} else {
				return []keyValuePair{}, errors.New("unsupported script")
			}
		}
		return pairs, nil
	} else if text, ok := commands.(string); ok {
		fmt.Printf("text: %s", text)
		// return Whens(strings.Split(text, ":")), nil
	} else if commands != nil {
		return []keyValuePair{}, errors.New("unsupported script")
	}

	return []keyValuePair{}, nil
}

func (c *JenkinsCiYamlParser) ParseYaml() (jenkinsCI *JenkinsCI, err error) {
	err = c.parseContent()
	if err != nil {
		return nil, err
	}

	jenkinsCI = &JenkinsCI{}

	err = c.loadOverall(jenkinsCI)

	if err != nil {
		return nil, err
	}

	err = c.loadPipeline()
	if err != nil {
		return nil, err
	}

	parsers := []struct {
		method func(pipeline *Pipeline) error
	}{
		{c.preparePipelineAgent},
		{c.preparePipelineEnv},
		{c.preparePipelineStages},
		//{c.prepareImage},
		//{c.prepareServices},
		//{c.prepareArtifacts},
		//{c.prepareCache},
	}
	pipeline := &Pipeline{}
	for _, parser := range parsers {
		err = parser.method(pipeline)
		if err != nil {
			return nil, err
		}
	}

	jenkinsCI.Pipeline = pipeline

	return jenkinsCI, nil
}

func NewGitLabCiYamlParser(data []byte) *JenkinsCiYamlParser {
	return &JenkinsCiYamlParser{
		data: data,
	}
}

type keyValuePair struct {
	Key   string
	Value string
}

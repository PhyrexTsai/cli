package main

import (
	"errors"
	"strings"
)

var HelpErr = errors.New("help")

type Cli struct {
	Topics   []*Topic
	Commands []*Command
}

func (cli *Cli) Parse(args []string) (ctx *Context, err error) {
	ctx = &Context{}
	if len(args) == 0 {
		return ctx, HelpErr
	}
	ctx.Topic, ctx.Command = cli.parseCmd(args[0])
	if ctx.Command == nil {
		return ctx, HelpErr
	}
	ctx.Args, ctx.App, err = parseArgs(ctx.Command, args[1:])
	return ctx, err
}

func (cli *Cli) parseCmd(cmd string) (topic *Topic, command *Command) {
	tc := strings.SplitN(cmd, ":", 2)
	topic = getTopic(tc[0], cli.Topics)
	if topic == nil {
		return nil, nil
	}
	if len(tc) == 2 {
		return topic, getCommand(tc[0], tc[1], cli.Commands)
	} else {
		return topic, getCommand(tc[0], "", cli.Commands)
	}
}

func getTopic(name string, topics []*Topic) *Topic {
	for _, topic := range topics {
		if topic.Name == name {
			return topic
		}
	}
	return nil
}

func getCommand(topic, command string, commands []*Command) *Command {
	for _, c := range commands {
		if c.Topic == topic && c.Command == command {
			return c
		}
	}
	return nil
}

func (cli *Cli) commandsForTopic(topic string) []*Command {
	commands := make([]*Command, 0, len(cli.Commands))
	for _, c := range cli.Commands {
		if c.Topic == topic {
			commands = append(commands, c)
		}
	}
	return commands
}

func parseArgs(command *Command, args []string) (result map[string]string, appName string, err error) {
	result = map[string]string{}
	numArgs := 0
	parseFlags := true
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "help" || args[i] == "--help" || args[i] == "-h":
			return nil, "", HelpErr
		case args[i] == "--":
			parseFlags = false
		case args[i] == "-a" || args[i] == "--app":
			i++
			if len(args) == i {
				return nil, "", errors.New("Must specify app name")
			}
			appName = args[i]
		case parseFlags && strings.HasPrefix(args[i], "-"):
			// TODO
		case numArgs == len(command.Args):
			return nil, "", errors.New("Unexpected argument: " + strings.Join(args[numArgs:], " "))
		default:
			result[command.Args[i].Name] = args[i]
			numArgs++
		}
	}
	for _, arg := range command.Args {
		if !arg.Optional && result[arg.Name] == "" {
			return nil, "", errors.New("Missing argument: " + strings.ToUpper(arg.Name))
		}
	}
	return result, appName, nil
}

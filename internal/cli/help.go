package cli

// Template for CLI help output
const HelpTemplate = `{{- if eq .CommandName "nodeadm"}}nodeadm{{else}}nodeadm {{.CommandName}}{{end}}{{if .Description}} - {{.Description}}
{{end}}

{{- if .PrependMessage}}
{{.PrependMessage}}
{{end}}

{{- if .UsageString}}
Usage:
  {{if ne .CommandName "nodeadm"}}nodeadm {{end}}{{.UsageString}}
{{end}}

{{- if .Positionals}}
Positional Variables:
{{- range .Positionals}}
  {{.Name}}  {{.Spacer}}{{if .Description}} {{.Description}}{{end}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{else}}{{if .Required}} (Required){{end}}{{end}}
{{- end}}
{{end}}

{{- if .Subcommands}}
Subcommands:
{{- range .Subcommands}}
  {{.LongName}}{{if .ShortName}} ({{.ShortName}}){{end}}{{if .Position}}{{if gt .Position 1}}  (position {{.Position}}){{end}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{end}}
{{- end}}
{{end}}

{{- if (gt (len .Flags) 0)}}
Flags:
{{- range .Flags}}
  {{if .ShortName}}-{{.ShortName}} {{else}}   {{end}}{{if .LongName}}--{{.LongName}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}{{end}}
{{- end}}
{{end}}

{{- if .AppendMessage}}
{{.AppendMessage}}
{{end}}

{{- if .Message}}
{{.Message}}
{{- end}}`




// const HelpTemplate = `{{if eq .CommandName "nodeadm"}}nodeadm{{else}}nodeadm {{.CommandName}}{{end}}{{if .Description}} - {{.Description}}{{end}}{{if .PrependMessage}}
// {{.PrependMessage}}{{end}}
// {{if .UsageString}}
// Usage:
//   {{if ne .CommandName "nodeadm"}}nodeadm {{end}}{{.UsageString}}{{end}}{{if .Positionals}}

// Positional Variables:{{range .Positionals}}
//   {{.Name}}  {{.Spacer}}{{if .Description}} {{.Description}}{{end}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{else}}{{if .Required}} (Required){{end}}{{end}}{{end}}{{end}}{{if .Subcommands}}

// Subcommands:{{range .Subcommands}}
//   {{.LongName}}{{if .ShortName}} ({{.ShortName}}){{end}}{{if .Position}}{{if gt .Position 1}}  (position {{.Position}}){{end}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{end}}{{end}}{{end}}{{if (gt (len .Flags) 0)}}

// Flags:{{range .Flags}}
//   {{if .ShortName}}-{{.ShortName}} {{else}}   {{end}}{{if .LongName}}--{{.LongName}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}{{end}}{{end}}{{end}}
// {{if .AppendMessage}}
// {{.AppendMessage}}{{end}}{{if .Message}}

// {{.Message}}{{end}}`


// const HelpTemplate = `{{if and .CommandName (ne .CommandName "nodeadm")}}{{.CommandName}}{{if .Description}} - {{.Description}}{{end}}{{else}}nodeadm{{if .Description}} - {{.Description}}{{end}}{{end}}{{if and .CommandName (ne .CommandName "nodeadm")}}{{.CommandName}}{{if .Description}} - {{.Description}}{{end}}{{else}}nodeadm{{if .Description}} - {{.Description}}{{end}}{{end}}{{if .PrependMessage}}
// {{.PrependMessage}}{{end}}
// {{if .UsageString}}
// Usage:
//   nodeadm {{.UsageString}}{{end}}{{if .Positionals}}

// Positional Variables:{{range .Positionals}}
//   {{.Name}}  {{.Spacer}}{{if .Description}} {{.Description}}{{end}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{else}}{{if .Required}} (Required){{end}}{{end}}{{end}}{{end}}{{if .Subcommands}}

// Subcommands:{{range .Subcommands}}
//   {{.LongName}}{{if .ShortName}} ({{.ShortName}}){{end}}{{if .Position}}{{if gt .Position 1}}  (position {{.Position}}){{end}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{end}}{{end}}{{end}}{{if (gt (len .Flags) 0)}}

// Flags:{{range .Flags}}
//   {{if .ShortName}}-{{.ShortName}} {{else}}   {{end}}{{if .LongName}}--{{.LongName}}{{end}}{{if .Description}}   {{.Spacer}}{{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}{{end}}{{end}}{{end}}
// {{if .AppendMessage}}
// {{.AppendMessage}}{{end}}{{if .Message}}

// {{.Message}}{{end}}`




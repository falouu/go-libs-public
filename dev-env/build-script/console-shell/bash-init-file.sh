source ~/.bashrc
{{if .AddToPath -}}
export PATH="{{.AddToPath}}:$PATH"
{{end -}}
PS1="[{{.PromptText}}] ${PS1}"

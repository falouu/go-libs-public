source ~/.bashrc
{{if .AddToPath -}}
export PATH="{{.AddToPath}}:$PATH"
{{end -}}
PS1="[{{.PromptText}}] ${PS1}"
export PS1_CUSTOM="{{.PromptText}}"   # to use with prompt engines (like starship.rs), that typically override PS1

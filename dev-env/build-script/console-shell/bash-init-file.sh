source ~/.bashrc
{{if .AddToPath -}}
export PATH="{{.AddToPath}}:$PATH"
{{end -}}
PS1="[buildscript] ${PS1}"
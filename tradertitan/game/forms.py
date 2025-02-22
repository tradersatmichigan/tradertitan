from django import forms

class JoinForm(forms.Form):
    name = forms.CharField(required = True, label="Username")

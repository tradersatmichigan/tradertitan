from django.contrib.auth import login
from django.http import HttpRequest, HttpResponse, JsonResponse
from django.shortcuts import render

from django.contrib.auth.models import User

from .forms import JoinForm

def signup(request: HttpRequest) -> HttpResponse:
    if request.method == "POST":
        form = JoinForm(request.POST)

        if form.is_valid():
            username = form.cleaned_data['name']

            if not User.objects.filter(username=username).exists():
                user = User.objects.create(username=username)
                login(request, user)
                return JsonResponse({"user": username, "id" : user.pk})
    else:
        form = JoinForm()
    return render(request, "signup.html", {"form" : form})
    

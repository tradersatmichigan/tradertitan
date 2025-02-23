from django.db import models
from django.contrib.auth.models import User

class Player(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    positionSum = models.IntegerField()
    totalPnl = models.IntegerField()


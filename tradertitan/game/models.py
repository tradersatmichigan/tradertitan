from django.db import models
from django.contrib.auth.models import User

class Player(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    room = models.ForeignKey("Room", on_delete=models.CASCADE)

    totalRank = models.IntegerField()
    totalPnl = models.IntegerField()


class Quote(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    price = models.IntegerField()

class Room(models.Model):
    bid = Quote()
    ask = Quote()

class Round(models.Model):
    market = models.CharField() # asset begin traded
    trueValue = models.IntegerField()

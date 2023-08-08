from fitness_pb2 import Fitness
import matplotlib.pyplot as plt

fitness = Fitness()

f = open("../data/fitness.bin", "rb")
fitness.ParseFromString(f.read())

plt.plot(fitness.AverageFitness)
plt.show()
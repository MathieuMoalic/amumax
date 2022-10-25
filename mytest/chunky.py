from llyr import op
from matplotlib import pyplot as plt

m = op("mytest/chunky.zarr")
m.p
arr = m.m[:]
plt.figure()
plt.imshow(arr[0, 0, :140, :, 1])
plt.show()

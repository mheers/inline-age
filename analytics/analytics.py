import matplotlib.pyplot as plt
import numpy as np
from mpl_toolkits.mplot3d import Axes3D

# Your dataset
data = [
    [2,1,292],
    [4,1,296],
    [8,1,300],
    [16,1,312],
    [32,1,332],
    [64,1,376],
    [2,2,684],
    [4,2,688],
    [8,2,696],
    [16,2,708],
    [32,2,732],
    [64,2,764],
    [2,3,1044],
    [4,3,1048],
    [8,3,1052],
    [16,3,1080],
    [32,3,1100],
    [64,3,1144],
    [2,4,1608],
    [4,4,1600],
    [8,4,1600],
    [16,4,1616],
    [32,4,1628],
    [64,4,1676],
    [2,5,1944],
    [4,5,1952],
    [8,5,1956],
    [16,5,1968],
    [32,5,1980],
    [64,5,2016],
    [2,6,2308],
    [4,6,2332],
    [8,6,2324],
    [16,6,2336],
    [32,6,2364],
    [64,6,2380],
    [2,7,2656],
    [4,7,2676],
    [8,7,2664],
    [16,7,2692],
    [32,7,2704],
    [64,7,2744],
    [2,8,3020],
    [4,8,3020],
    [8,8,3032],
    [16,8,3064],
    [32,8,3080],
    [64,8,3096],
    [2,9,3736],
    [4,9,3752],
    [8,9,3752],
    [16,9,3776],
    [32,9,3760],
    [64,9,3812],
    [2,10,4092],
    [4,10,4100],
    [8,10,4104],
    [16,10,4112],
    [32,10,4132],
    [64,10,4176],
    [2,11,4640],
    [4,11,4620],
    [8,11,4620],
    [16,11,4652],
    [32,11,4652],
    [64,11,4708],
    [2,12,5176],
    [4,12,5172],
    [8,12,5176],
    [16,12,5184],
    [32,12,5204],
    [64,12,5240],
    [2,13,5700],
    [4,13,5692],
    [8,13,5728],
    [16,13,5716],
    [32,13,5740],
    [64,13,5780],
    [2,14,6232],
    [4,14,6220],
    [8,14,6228],
    [16,14,6240],
    [32,14,6248],
    [64,14,6316],
    [2,15,6760],
    [4,15,6756],
    [8,15,6768],
    [16,15,6772],
    [32,15,6804],
    [64,15,6848],
]

# Separate the data into three lists for m, p, and e
m_values, p_values, e_values = zip(*data)

# Convert the lists to NumPy arrays
m_values = np.array(m_values)
p_values = np.array(p_values)
e_values = np.array(e_values)

# Create a 3D plot
fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

# Plot the data
ax.scatter(m_values, p_values, e_values, c='r', marker='o')

# Set labels for each axis
ax.set_xlabel('m (Length of the Message)')
ax.set_ylabel('p (Number of Participants)')
ax.set_zlabel('e (Length of Encrypted Message)')

# Show the plot
plt.show()
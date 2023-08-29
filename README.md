500 Watt Antenna Tuner - Simulator
The simulator has three functions:

Calculate the exact inductor and capacitor values for the matching network using the analytical model.  It also calculates the maximum values of inductors and capacitors for the top and bottom of each amateur HF (and 6 meter) band.
Approximate the exact inductor and capacitor values using the power of two weighted eight standard inductor and capacitor values.
Calculates the voltage across the capacitors and the current through the inductors.  At the same time, it calculates the current through each specific standard capacitor and the voltage across each standard inductor.
The simulator is written in Go programming language.  You can find the code on Github:

https://github.com/Saied74/tuner

If you know C/C++, the code should be easy to read.  You should also consider learning Go.  It is compiled, strongly typed, object oriented language.  It is compact, expressive, and functional with a rich set of libraries.  

The simulator depends on my command line library that can be found at:

https://github.com/Saied74/cli

I recently noticed that Go has a command line library.  I will have to look into it and maybe switch at some point.

I will describe the operation of simulator through its user interface.  It has a number of commands that it displays in a numbered list.  By entering the number, you will get a prompt to enter one of the options in that feature or command.  It validates the option and if is not a valid option it asks you to enter it again.

Here is the list of the available commands:

1. Change the csv file name = data.csv
2. Change the minMax file name = minMax.csv
3. Source voltage in volts = 223
4. Capacitor Q = 1000
5. Inductor Q = 100
6. Calculate simple LC and LC min max calculation
7. Fit LC to their standard values and calculate max difference to actual values
8. Calculate current and voltage across/through L, C, and relays and their maximum values
Each line shows the command that it invokes and also the default value for that command after the = sign.

The last three lines (6, 7, and 8) execute the three scenarios listed on top of this page (exact L & C calculations, approximation of the L & C values using the standard L and C ladder and currents through and voltages across these components.  The result is written into the file named in line 1.  In each case, the minimum and maximum values of these items are also calculated and the result is written into the file named in line 2.  These files are stored in the directory which is hard coded:

 ~/Documents/hamradio/Antennas/tuner/Simulation_output

The standard L and C ladder values are also hard coded.  If you would like to change either one of these, unfortunately, you will have to edit the source.

Line 3 is the peak sinusoidal voltage applied to the input of the tuner.  For a 500 watt output into 50 ohms, it would be 223 volts.  

The data file is one header line and 3600 rows as the simulator steps through 10 values of SWR and 360 values of the reflection coefficient phase angle.  The max file contains the maximum values of inductors and capacitors or voltages and currents.

The column headings for the first file are:

swr: value of the SWR
theta: reflection coefficient phase angle
r0: normalized (to 50 ohms) load resistance
x0: normalized (to 50 ohms) load reactance
r1: resistance coordinate of intersection with the r = 1 or g = 1 circles
x1: reactance coordinate of intersection with the r = 1 or g = 1 circles
region: the region where load is located as described in the analytical model
parallel: normalized (to 50 ohms) parallel reactance required for match
series: normalized (to 50 ohms) series reactance required for match
Columns J through O are not used
Columns P through AI the inductance and capacitance to achieve a match at the low and high end of each amateur band.  This is the case for both simple and fit LC calculations.
For VI calculations, columns A through I are the same as above.  But columns P through OU show the voltage and currents through inductors and capacitors in each band.  For each end of the band, the capacitor voltage followed by all the capacitor currents are shown in successive columns.  This is followed by inductor current, followed by successive columns of inductor voltages.  The capacitor voltage and inductor current columns are so noted in the headings. 
The minMax file format is much simpler and self explanatory.  

This is a short summary of how to use the simulator.  You can of course read through the source code and modify and enhance it as you need.  I will describe the structure of the code in the comments to the code.


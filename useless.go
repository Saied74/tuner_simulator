package main

// "oneError": &cli.Item{
//   Name:      "oneError",
//   Prompt:    "Run the simulation with one errors and write a csv file",
//   Response:  "Do I need this 1?",
//   Value:     "",
//   Validator: cli.ItemValidator(func(x string) bool { return true }),
// },

// case "oneError":
//   f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
//   if err != nil {
//     log.Fatal(err)
//   }
//   defer f.Close()
//   err = writeImpedanceHeader(f)
//   if err != nil {
//     log.Fatal(err)
//   }
//   _, err = f.WriteString("swr")
//   if err != nil {
//     log.Fatal(err)
//   }
//   _, err = f.WriteString("\n")
//   if err != nil {
//     log.Fatal(err)
//   }
//   for _, w := range swr {
//     s.s = w
//     gamma := (s.s - 1.0) / (s.s + 1.0)
//     s.gamma = gamma
//     if s.which == "gamma" {
//       s.gamma += s.gainTol * gamma
//     }
//     for i := 0; i < 360; i++ {
//       theta := float64(i)
//       s.theta = theta
//       if s.which == "theta" {
//         s.theta += s.phaseTol
//       }
//       s.trueCalc()
//       err = s.writeImpedance(f)
//       if err != nil {
//         log.Fatal(err)
//       }
//       swr := calcSWR(s.point1.r, s.point1.x)
//       _, err = f.WriteString(fmt.Sprintf("%.2f", swr))
//       if err != nil {
//         log.Fatal(err)
//       }
//       _, err = f.WriteString("\n")
//       if err != nil {
//         log.Fatal(err)
//       }
//     }
//   }

// "distance": &cli.Item{
//   Name:      "distance",
//   Prompt:    "Calculate the minimum and maximum distances",
//   Response:  "Do I need this 3?",
//   Value:     "",
//   Validator: cli.ItemValidator(func(x string) bool { return true }),
// },

// case "distance":
//   f, err := os.OpenFile(s.outputFile, os.O_RDWR|os.O_CREATE, 0666)
//   if err != nil {
//     log.Fatal(err)
//   }
//   defer f.Close()
//
//   err = writeDistanceHeader(f)
//   if err != nil {
//     log.Fatal(err)
//   }
//   _, err = f.WriteString("\n")
//   if err != nil {
//     log.Fatal(err)
//   }
//
//   for _, w := range swr {
//     s.s = w
//     s.gamma = (s.s - 1.0) / (s.s + 1.0)
//     for i := 0; i < 360; i++ {
//       s.theta = float64(i)
//       s.trueCalc()
//       switch s.region {
//       case 1:
//         if s.seriesReact > s.baseMaxSeries1.seriesReact {
//           s.copyExt(s.baseMaxSeries1)
//         }
//         if s.seriesReact < s.baseMinSeries1.seriesReact {
//           s.copyExt(s.baseMinSeries1)
//         }
//         if s.parallelReact > s.baseMaxParallel1.parallelReact {
//           s.copyExt(s.baseMaxParallel1)
//         }
//         if s.parallelReact < s.baseMinParallel1.parallelReact {
//           s.copyExt(s.baseMinParallel1)
//         }
//       case 2:
//         if s.seriesReact > s.baseMaxSeries2.seriesReact {
//           s.copyExt(s.baseMaxSeries2)
//         }
//         if s.seriesReact < s.baseMinSeries2.seriesReact {
//           s.copyExt(s.baseMinSeries2)
//         }
//         if s.parallelReact > s.baseMaxParallel2.parallelReact {
//           s.copyExt(s.baseMaxParallel2)
//         }
//         if s.parallelReact < s.baseMinParallel2.parallelReact {
//           s.copyExt(s.baseMinParallel2)
//         }
//       }
//       s.gammaTemp = s.gamma
//       s.thetaTemp = s.theta
//       switch s.which {
//       case "theta":
//         s.theta += s.phaseTol
//       case "gamma":
//         s.gamma += s.gamma * s.gainTol
//       }
//       s.trueCalc()
//       switch s.region {
//       case 1:
//         if s.seriesReact > s.tolMaxSeries1.seriesReact {
//           s.copyExt(s.tolMaxSeries1)
//         }
//         if s.seriesReact < s.tolMinSeries1.seriesReact {
//           s.copyExt(s.tolMinSeries1)
//         }
//         if s.parallelReact > s.tolMaxParallel1.parallelReact {
//           s.copyExt(s.tolMaxParallel1)
//         }
//         if s.parallelReact < s.tolMinParallel1.parallelReact {
//           s.copyExt(s.tolMinParallel1)
//         }
//       case 2:
//         if s.seriesReact > s.tolMaxSeries2.seriesReact {
//           s.copyExt(s.tolMaxSeries2)
//         }
//         if s.seriesReact < s.tolMinSeries2.seriesReact {
//           s.copyExt(s.tolMinSeries2)
//         }
//         if s.parallelReact > s.tolMaxParallel2.parallelReact {
//           s.copyExt(s.tolMaxParallel2)
//         }
//         if s.parallelReact < s.tolMinParallel2.parallelReact {
//           s.copyExt(s.tolMinParallel2)
//         }
//       }
//       s.gamma = s.gammaTemp
//       s.theta = s.thetaTemp
//     }
//   }
//   err = s.writeDistance(f)
//   if err != nil {
//     log.Fatal(err)
//   }
//
// "iterations": &cli.Item{
//   Name:      "iterations",
//   Prompt:    "How many iterations for swr calculaton",
//   Response:  "Do I need this 4?",
//   Value:     "2",
//   Validator: iterValidator,
// },
//
// var iterValidator = cli.ItemValidator(func(x string) bool {
// 	y, err := strconv.Atoi(x)
// 	if err != nil {
// 		return false
// 	}
// 	if y < 1 {
// 		return false
// 	}
// 	return true
// })
//
// case "iterations":
//   iter, _ := strconv.Atoi(item.Value)
//   s.iteration = iter
//   fmt.Println("Iteration: ", s.iteration)

#This is the server file for the coperformance project
library(shiny)
library(shinyjs)

function(input, output, session) {
	rvals = reactiveValues()
	rvals$RUNPLOT = FALSE
	rvals$jsonname = "_"
	rvals$cperfname = "_"
	rvals$performanceList = list()
	
	# Disable the data download by default
	shinyjs::disable("djson")
	shinyjs::disable("dcoperf")
	
	# Define the input and output files
	inputFile = "complete.json"
	outputFile = paste(BaseDir, "CompJSON.json", sep="/")
	
	observeEvent(input$FromSeason, handlerExpr = {
		updateSelectizeInput(session, "ToSeason", 
			choices=yvec[which(yvec == input$FromSeason):length(yvec)], selected=input$ToSeason
		)
	})
	
	observeEvent(input$makePlot, handlerExpr = {
	  # To start, define the inputs to the Go function
	  fileInfo = c(inputFile, outputFile)
	  minPrograms = as.integer(input$MinPerf)
	  seasonRange = paste(input$FromSeason, input$ToSeason, sep=":")
	  
	  # Now run it
	  goOutput = .Call("ProcessNYPhilJSON", fileInfo, minPrograms, seasonRange)
	  if (goOutput[[4]] != "nil") {
	    output$ResMessage = renderText(paste("ERROR:", goOutput[[4]]))
	    return()
	  }
	  
	  # Now make the network itself
	  networkPlot = MakeNetwork(goOutput)
	  
	  # Now send the plot on its way
	  output$network = renderVisNetwork(networkPlot)
	  output$ResMessage = renderText("No errors")
	  
	  #Set the reactive values for other functions
	  rvals$RUNPLOT = TRUE
	  rvals$Cframe = goOutput[[1]]
	  rvals$cperfname = paste(paste("Coperformance", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
	                          ".csv", sep=""
	  )
	  rvals$jsonname = paste(paste("NYPhilParsed", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
	                         ".json", sep=""
	  )
	  rvals$performanceList = goOutput
	  
	  print(goOutput[[5]])
	  if (goOutput[[5]] == "nil") {
	    shinyjs::enable("djson")
	    shinyjs::enable("dcoperf")
	  } else {
	    shinyjs::disable("djson")
	    shinyjs::disable("dcoperf")
	  }
	  
	})
	
	# Handle the download of the json file
	output$djson = downloadHandler(
	  filename = function() {rvals$jsonname},
	  content = function(file) {
	      jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
	      file.copy(jsonfile, file)
	  }
	)
	
	# Handle the csv download of the coperformance matrix
	output$dcoperf = downloadHandler(
	  filename = function() {rvals$cperfname},
	  content = function(file) {
	    tempFile = paste(BaseDir, "coperf.csv", sep="/")
	    
	    cMatrix = matrix(rvals$performanceList[[1]][[2]], 
	                     nrow=rvals$performanceList[[1]][[1]][1], 
	                     ncol=rvals$performanceList[[1]][[1]][2]
      )
	    csvFrame = data.frame(cMatrix)
	    names(csvFrame) = rvals$performanceList[[2]]
	    row.names(csvFrame) = rvals$performanceList[[2]]

	    write.csv(csvFrame, tempFile, row.names=T)
	    file.copy(tempFile, file)
	  }
	)
	
}
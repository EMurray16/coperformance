#This is the server file for the coperformance project
library(shiny); library(visNetwork)

#Load the extras
dyn.load("processPhil.so")
yvec = read.csv("philyears.csv", stringsAsFactors=F, header=F)[,1]
source("GoHelp.r")
source("MakeNetwork.r")
BaseDir = getwd()

function(input, output, session) {
	observeEvent(input$FromSeason, handlerExpr = {
		updateSelectizeInput(session, "ToSeason", 
			choices=yvec[which(yvec == input$FromSeason):length(yvec)], selected=input$ToSeason
		)
	})
	
	observeEvent(input$makePlot, handlerExpr = {
		#Set the inputs to the workhorse
		yearstring = paste(input$FromSeason, input$ToSeason, sep=":")
		MinPerf = as.integer(input$MinPerf)
		InputFile = paste(BaseDir, "complete.json", sep="/")
		jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
		csvfile = paste(BaseDir, "Coperf.csv", sep="/")
		
		#Encode the inputs
		InYear = GoEncode(yearstring)
		InFile = GoEncode(InputFile)
		JFile = GoEncode(jsonfile)
		CFile = GoEncode(csvfile)
		
		#Now call the Go function
		res = .Call("ProcessPhil", MinPerf, InYear, InFile, JFile, CFile)
		
		#Handle any errors
		errstring = GoDecode(res)
		errs = strsplit(errstring, " /AND/ ")[[1]]
		for (i in 1:3) {
			if (errs[i] != "nil") {
				output$ResMessage = renderText("Error loading and analyzing data, see below:")
				output$ResPrint = renderPrint(errs[i])
				return()
			}
		}
		
		#Now make the network
		NetList = MakeNetwork(csvfile)
		if (NetList[[2]] >= 7) {
			output$NetMessage = renderText("Too many clusters to visualize effectively, showing vanilla network.")
		} else {
			sentence = paste("There are", NetList[[2]], "clusters in the network.", sep=" ") 
			output$NetMessage = renderText(paste("There are", NetList[[2]], "clusters in the network.", sep=" "))
		}
		output$network = renderVisNetwork(NetList[[1]])
	})
	
	output$djson = downloadHandler(
		filename = paste(
			paste("NYPhilParsed", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
			".json", sep=""
		),
		content = function(file) {
			#Set the inputs to the workhorse
			yearstring = paste(input$FromSeason, input$ToSeason, sep=":")
			MinPerf = as.integer(input$MinPerf)
			InputFile = paste(BaseDir, "complete.json", sep="/")
			jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
			csvfile = paste(BaseDir, "Coperf.csv", sep="/")
			
			#Encode the inputs
			InYear = GoEncode(yearstring)
			InFile = GoEncode(InputFile)
			JFile = GoEncode(jsonfile)
			CFile = GoEncode(csvfile)
			
			#Now call the Go function
			res = .Call("ProcessPhil", MinPerf, InYear, InFile, JFile, CFile)
			
			#Handle any errors
			errstring = GoDecode(res)
			errs = strsplit(errstring, " /AND/ ")[[1]]
			for (i in 1:3) {
				if (errs[i] != "nil") {
					output$ResMessage = renderText("Error loading and analyzing data, see below:")
					output$ResPrint = renderPrint(errs[i])
					return()
				}
			}
			
			file.copy("CompJSON.json", file)
		}
	)
	
	output$dcoperf = downloadHandler(
		filename = paste(
			paste("Coperformance", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
			".csv", sep=""
		),
		content = function(file) {
			#Set the inputs to the workhorse
			yearstring = paste(input$FromSeason, input$ToSeason, sep=":")
			MinPerf = as.integer(input$MinPerf)
			InputFile = paste(BaseDir, "complete.json", sep="/")
			jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
			csvfile = paste(BaseDir, "Coperf.csv", sep="/")
			
			#Encode the inputs
			InYear = GoEncode(yearstring)
			InFile = GoEncode(InputFile)
			JFile = GoEncode(jsonfile)
			CFile = GoEncode(csvfile)
			
			#Now call the Go function
			res = .Call("ProcessPhil", MinPerf, InYear, InFile, JFile, CFile)
			
			#Handle any errors
			errstring = GoDecode(res)
			errs = strsplit(errstring, " /AND/ ")[[1]]
			for (i in 1:3) {
				if (errs[i] != "nil") {
					output$ResMessage = renderText("Error loading and analyzing data, see below:")
					output$ResPrint = renderPrint(errs[i])
					return()
				}
			}
			
			file.copy("Coperf.csv", file)
		}
	)
	
}
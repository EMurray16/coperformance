#This is the server file for the coperformance project
library(shiny); library(visNetwork)

#Load the extras
dyn.load("processPhil2.so")
yvec = read.csv("philyears.csv", stringsAsFactors=F, header=F)[,1]
source("GoHelp.r")
source("MakeNetwork.r")
BaseDir = "$USER/ShinyApps/coperformance/temp"
HomeDir = getwd()

function(input, output, session) {
	rvals = reactiveValues()
	rvals$RUNPLOT = FALSE
	rvals$jsonname = "_"
	rvals$cperfname = "_"
	
	observeEvent(input$FromSeason, handlerExpr = {
		updateSelectizeInput(session, "ToSeason", 
			choices=yvec[which(yvec == input$FromSeason):length(yvec)], selected=input$ToSeason
		)
	})
	
	observeEvent(input$makePlot, handlerExpr = {
		#Set the inputs to the workhorse
		yearstring = paste(input$FromSeason, input$ToSeason, sep=":")
		MinPerf = as.integer(input$MinPerf)
		InputFile = paste(HomeDir, "complete.json", sep="/")
		jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
		
		#Encode the inputs
		InYear = GoEncode(yearstring)
		InFile = GoEncode(InputFile)
		JFile = GoEncode(jsonfile)
		
		#Now call the Go function to get the coperformacne matrix
		res = .Call("GetCoperformance", MinPerf, InYear, InFile, JFile)
		
		#Handle any errors
		errstring = GoDecode(res[[3]])
		prefix = GoDecode(res[[3]][1:8])
		if (prefix == "WARNING:") {
			output$ResMessage = renderText("Warning loading and analyzing data:")
			output$ResPrint = renderPrint(errstring)
		} else if (errstring != "nil") {
			output$ResMessage = renderText("Error loading and analyzing data, see below:")
			output$ResPrint = renderPrint(errstring)
			return()
		}
		
		#Make the coperformance data frame
		Cframe = data.frame(Matricize(res[[1]]))
		Cnames = String2vec(GoDecode(res[[2]]), delim="~")
		names(Cframe) = Cnames; row.names(Cframe) = Cnames
		
		#Now make the network
		NetList = MakeNetwork(Cframe)
		if (NetList[[2]] >= 7) {
			output$NetMessage = renderText("Too many clusters to visualize effectively, showing vanilla network.")
		} else {
			sentence = paste("There are", NetList[[2]], "clusters in the network.", sep=" ") 
			output$NetMessage = renderText(paste("There are", NetList[[2]], "clusters in the network.", sep=" "))
		}
		output$network = renderVisNetwork(NetList[[1]])
		
		#Set the reactive values for other functions
		rvals$RUNPLOT = TRUE
		rvals$Cframe = Cframe
		rvals$cperfname = paste(paste("Coperformance", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
			".csv", sep=""
		)
		rvals$jsonname = paste(paste("NYPhilParsed", input$MinPerf, input$FromSeason, input$ToSeason, sep="_"),
			".json", sep=""
		)
	})
	
	output$djson = downloadHandler(
		filename = function() {rvals$jsonname},
		content = function(file) {
			#Make sure the plot function has been run
			if (!rvals$RUNPLOT) {
				output$ResMessage = renderText("Main plotting function must be run before donwloading file!")
			} else {
				jsonfile = paste(BaseDir, "CompJSON.json", sep="/")
				file.copy(jsonfile, file)
			}
		}
	)
	
	output$dcoperf = downloadHandler(
		filename = function() {rvals$cperfname},
		content = function(file) {
			#Make sure the plot function has been run
			if (!rvals$RUNPLOT) {
				output$ResMessage = renderText("Main plotting function must be run before donwloading file!")
			} else {
				print(rvals$cperfname)
				csvfile = paste(BaseDir, "Coperf.csv", sep="/")
				#Write the csv file
				write.csv(rvals$Cframe, csvfile, row.names=TRUE)
				file.copy(csvfile, file)
			}
		}
	)
	
}
#This is for the user interface of the unnamed NY Philharmonic project
library(shiny); library(shinythemes); library(visNetwork)
yvec = read.csv("philyears.csv", stringsAsFactors=F, header=F)[,1]


shinyUI( fluidPage( 
	headerPanel(strong("New York Philharmonic Coperformance Project")),
	sidebarPanel(h4(strong("Data Selection")),
		p("Coperformance: The ratio between the number programs two composers appear in 
			together and the number of times they should appear together at random"
		),
		hr(),
		selectizeInput("FromSeason", label="From", choices=yvec, width="50%", multiple=F, selected="1842-43"),
		selectizeInput("ToSeason", label="To", choices=yvec, width="50%", multiple=F, selected="2017-18"),
		numericInput("MinPerf", label="Minimum Performances", min=1,max=1000,step=1,value=400),
		actionButton("makePlot", "Plot Network"),
		hr(),
		h4(strong("Download Data")),
		downloadButton("djson", label="Download json"),
		downloadButton("dcoperf", label="Download matrix"),
		hr(),
		p("For more information, see my ", a("Github repository", href="https://github.com/EMurray16/coperformance"), 
			" and my ", a("blog post", href="https://overthinkdciscores.com/2018/12/02/using-data-to-visualize-connections-between-composers/"),
			" on the project."
		)
	),
	mainPanel(
		textOutput("ResMessage"),
		verbatimTextOutput("ResPrint"),
		textOutput("NetMessage"),
		visNetworkOutput("network")
	),
	title="NYPhil_SNA", theme=shinytheme("yeti")
))
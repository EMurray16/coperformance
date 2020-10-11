# This defines the global environment for the shiny app
libs = c("data.table","igraph","visNetwork")
lapply(libs, library, character.only=T)
# Load the dynamic library
if (!is.loaded("ProcessNYPhilJSON")) {
  dyn.load("nyphilLibrary.so")
}

yvec = read.csv("philyears.csv", stringsAsFactors=F, header=F)[,1]

BaseDir = getwd()
HomeDir = getwd()

ColorList = list(rp=rgb(0.8,0.6,0.7), vm=rgb(0.8,0.4,0), bg=rgb(0,0.6,0.5), yl=rgb(0.95,0.9,0.25), bl=rgb(0,0.45,0.7), 
               or=rgb(0.9,0.6,0), sb=rgb(0.35,0.7,0.9)
)

# MakeNetwork creates a network plot based on the Go output
MakeNetwork <- function(goOutput) {
  composerNames = goOutput[[2]]
  networkFrom = goOutput[[3]][[1]]
  networkTo = goOutput[[3]][[2]]
  connectionStrength = goOutput[[3]][[3]]
  composerCount = goOutput[[3]][[4]]
  
  composerNodes = data.table(id=composerNames, value=composerCount, title=composerNames)
  composerEdges = data.table(from=networkFrom, to=networkTo, width=connectionStrength*2, color.opacity=connectionStrength, color="black")
  
  #Make an Igraph object for clustering
  CompIgraph = graph_from_data_frame(d=composerEdges, vertices=composerNodes, directed=F)
  CompClust = cluster_fast_greedy(CompIgraph, weights=connectionStrength)
  #Extract the clusters
  ClustSum = communities(CompClust)
  
  #Set the colors of the clusters to each composers
  if (length(ClustSum) <= 7) { #we can set the colors
    composerNodes$color = rep(0, length(composerNames))
    for (i in 1:length(ClustSum)) {
      composerNodes$color[which(composerNodes$id %in% ClustSum[[i]])] = ColorList[[i]]
    }
  }
  
  cnet = visNetwork(composerNodes, composerEdges, main="Composer Coperformance Network") %>%
    visIgraphLayout(layout="layout_with_lgl") %>%
    visOptions(highlightNearest=T, selectedBy="id")
  return(cnet)
}
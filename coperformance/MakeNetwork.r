#This contains the function to make the network given the coperformance file
Libraries = c("visNetwork","igraph")
lapply(Libraries, library, character.only=T)
ColList = list(rp=rgb(0.8,0.6,0.7), vm=rgb(0.8,0.4,0), bg=rgb(0,0.6,0.5), yl=rgb(0.95,0.9,0.25), bl=rgb(0,0.45,0.7), 
	or=rgb(0.9,0.6,0), sb=rgb(0.35,0.7,0.9)
)

MakeNetwork <- function(CoperfFile) {
	CoperfMat = read.csv(CoperfFile, row.names=1)
	
	#Now make the matrix into a data frame
	nodes = data.frame(id=names(CoperfMat), value=rep(1,nrow(CoperfMat)), label=rep("",nrow(CoperfMat)))
	nodes$title = names(CoperfMat)
	#Mkae the edges
	IDs1 = vector(mode='character')
	IDs2 = vector(mode='character')
	Strengths = vector(mode='double')
	for (R in 2:nrow(CoperfMat)) {
		for (C in 1:R) {
			if (CoperfMat[R,C] > 1 & R != C) {
				IDs1 = c(IDs1, names(CoperfMat)[R])
				IDs2 = c(IDs2, names(CoperfMat)[C])
				Strengths = c(Strengths, CoperfMat[R,C])
			}
			if (C == R) {
				nodes[R,"value"] = CoperfMat[R,C]
			}
		}
	}
	edges = data.frame(from=IDs1, to=IDs2, width=Strengths, color="black")
	edges$color.opacity = Strengths

	#Make an Igraph object for clustering
	CompIgraph = graph_from_data_frame(d=edges, vertices=nodes, directed=F)
	CompClust = cluster_fast_greedy(CompIgraph, weights=Strengths)
	#Extract the clusters
	ClustSum = communities(CompClust)
	
	#Set the colors of the clusters to each composers
	if (length(ClustSum) <= 7) { #we can set the colors
		nodes$color = rep(0,nrow(CoperfMat))
		for (i in 1:length(ClustSum)) {
			nodes$color[which(nodes$id %in% ClustSum[[i]])] = ColList[[i]]
		}
	}

	cnet = visNetwork(nodes, edges, main="Composer Coperformance Network") %>%
		visIgraphLayout(layout="layout_with_lgl") %>%
		visOptions(highlightNearest=T, selectedBy="id")
	#print(cnet)
	return(list(cnet, length(ClustSum)))
}



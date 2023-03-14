// // When the button is clicked, make a request to the Go server and create a mind map with the response data
// d3.select("button").on("click", function() {
//     d3.json("http://localhost:8080/data", function(error, data) {
//         if (error) {
//             console.log(error);
//             return;
//         }
        
//         // Create the mind map
//         var mindmap = d3.select("#mind-map");
        
//         var nodes = mindmap.selectAll("div")
//             .data(data.nodes)
//             .enter()
//             .append("div")
//             .attr("class", "node")
//             .html(function(d) {
//                 return "<h4>" + d.id + "</h4>" +
//                        "<p>Ingressing: " + d.ingressing + "</p>" +
//                        "<p>Egressing: " + d.egressing + "</p>";
//             });
//     });
// });

request = document.querySelector(".request");

request.onclick = function() {
    this.innerHTML = "<div class='loader'> </div>";
    setTimeout(() => {
        this.innerHTML = "Successfull request"
        this.style = "background: green; color white; pointer-events: none;"
    }, 25000);
}


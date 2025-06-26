transduce-go

This is a transducer library based heavily on a [java library](https://github.com/cognitect-labs/transducers-java) for how to go about this. Due to Golang not having dynamic dispatch of function return types in interfaces composition of transducers have to specify return types limiting their reusability depending on the code base but can still be helpful in isolating business logic to funtions.

I did this mostly because of scenarios where I wanted to decompose logic into reusable generic pieces of code and to get around having to know the full types of the pipeline with functions chaining, transducers offer the best path despite being very confusing to grasp when trying to leverage them. So this doesn't save you headache of code composition, it just separates the business logic from the composition of different functions a bit better. It's also important to note that this currently does not support any concurrency as this was originally intended for handling sequential data (namely pagination).

For more info on transducers (what they are, how they work, and why bother) [this medium post](https://medium.com/javascript-scene/transducers-efficient-data-processing-pipelines-in-javascript-7985330fe73d) does a decent job setting the scene.

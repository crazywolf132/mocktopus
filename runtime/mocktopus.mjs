/**
* This is the middleware system that allows for mocktopus to run
* javascript. We will read and execute the route file, and return
* the response as JSON format to the console that the `mocktopus`
* application can then read.
*
* This allows for us to get around the usage of the nodeJS runtime.
*/
import { readFileSync } from 'node:fs';
import { } from 'node:path'

(() => {
  const fileToRun = process.argv[2];
  if (fileToRun === "") {
    throw new Error("No file provided")
  }

  // We know the location, so we will read in the contents.
  const contents = readFileSync(fileToRun, 'utf8');
  // We are going to create an exports overloader.
  const module = { exports: {}};
  // We are now going to execute the contents of the file.
  Function("module", "exports", contents)(module, module.exports);
  // Grabbing the results.
  const export = "default" in module?.exports ? module?.exports?.default : module.exports;
  // If the export is the default export and it is a function, we will execute it and record the result.
  let result;
  if (typeof export === "function") {
    result = export();
  }
  // We will now convert what we can into a JSON stringy and console log for it for the parent system to run.
  console.log(JSON.stringify(result));
})();

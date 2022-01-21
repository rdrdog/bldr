# Plugins

Plugins are specific types of jobs that are run as part of the build or deployment phase of the pipeline. Like Github actions, there are a few builtin ones, and we will support external repositories soon!

Unlike Github actions, plugins are run with a specific (opinionated) set of contextual details passed from previous steps in the pipeline. While this can be done using Github actions, there are limitations (file-based context passing between jobs, only simple strings can be used as inputs and outputs).

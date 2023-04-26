# Morty

![Morty banner](assets/morty.jpg)

Morty is an open source serverless platform for managing functions as a service. It is mainly written in Go, and it use under the hood the [RIK](https://github.com/rik-org/rik) orchestrator, to manage microVM instances where the functions will be executed.

## Features

- **Manage functions in a simple way**: With the official [Morty CLI](https://github.com/morty-faas/cli), you can create as many functions as you want. In 3 commands, you will be able to invoke your first function in a blazingly fast time.

- **Native support for various runtimes**: You like to code serverless functions in `NodeJS`, `Go`, `Rust` or `Python` ? Morty has native support for them. If you wish to develop functions in another runtime, you can open an issue in the [runtimes repository](https://github.com/morty-faas/runtimes).

- **Community driven**: Morty is an open-source project and contributions are highly appreciated. You need to ask a question, you want to make a proposal for a new feature or you want to fix a bug ? Don&apos;t hesitate to do it on GitHub ! [Learn more about our contribution process](#contributing).

## Documentation

You will find all the resources needed to start using Morty or learn the internal design in our official documentation avalaible on https://morty-faas.github.io.

The source code of the documentation is available on the [morty-faas/morty-faas.github.io](https://morty-faas.github.io)

## Contributing

We want to make Morty a community-driven platform where everyone should be able to contribute. No matter what type of contribution you want to bring to Morty, be sure that we will consider it.

You can learn more about our contribution guidelines in the following [document](https://morty-faas.github.io/contribute).

## License

Morty is [MIT licensed](./LICENSE).

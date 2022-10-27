// swift-tools-version:5.3
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "Suborbital",
    products: [
        .library(name: "Suborbital", targets: ["Suborbital"]),
    ],
    dependencies: [],
    targets: [
        .target(
            name: "Suborbital",
            dependencies: [],
            path: "sdk/swift/Sources"),
    ]
)

import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.0"
    application
}

group = "me.scmacdon"
version = "1.0-SNAPSHOT"

java {
    sourceCompatibility = JavaVersion.VERSION_17
    targetCompatibility = JavaVersion.VERSION_17
}

buildscript {
    repositories {
        maven("https://plugins.gradle.org/m2/")
    }
    dependencies {
        classpath("org.jlleitschuh.gradle:ktlint-gradle:10.3.0")
    }
}

repositories {
    mavenCentral()
}
apply(plugin = "org.jlleitschuh.gradle.ktlint")
dependencies {
    implementation("aws.sdk.kotlin:support:1.0.30")
    implementation("aws.smithy.kotlin:http-client-engine-okhttp:0.30.0")
    implementation("aws.smithy.kotlin:http-client-engine-crt:0.30.0")
    testImplementation("org.junit.jupiter:junit-jupiter:5.9.2")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.6.4")
}
tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}

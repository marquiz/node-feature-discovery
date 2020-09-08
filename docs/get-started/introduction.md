---
title: "Introduction"
layout: default
parent: Get started
nav_order: 1
---

# Introduction
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

This software enables node feature discovery for Kubernetes. It detects
hardware features available on each node in a Kubernetes cluster, and
advertises those features using node labels.

NFD consists of two software components:
1. nfd-master
2. nfd-worker


## NFD-Master

Nfd-master is the daemon responsible for communication towards the Kubernetes
API. That is, it receives labeling requests from the worker and modifies node
objects accordingly.

## NFD-Worker

Nfd-worker is a daemon responsible for feature detection. It then communicates
the information to nfd-master which does the actual node labeling.  One
instance of nfd-worker is supposed to be running on each node of the cluster,

## Feature Discovery

Quick introduction to feature sources and the labels.

## Node Annotations

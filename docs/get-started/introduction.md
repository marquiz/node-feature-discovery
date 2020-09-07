---
title: "Introduction"
layout: default
parent: Get started
nav_order: 2
---

# Introduction
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

This software enables node feature discovery for Kubernetes. It detects
hardware features available on each node in a Kubernetes cluster, and advertises
those features using node labels.

NFD consists of two software components:
1. **nfd-master** is responsible for labeling Kubernetes node objects
2. **nfd-worker** is detects features and communicates them to nfd-master.
   One instance of nfd-worker is supposed to be run on each node of the cluster


## NFD-Master


## NFD-Worker


## Feature Discovery

Quick introduction to feature sources and the labels.

## Node Annotations

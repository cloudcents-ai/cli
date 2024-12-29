# Cloud Cents CLI ![Go](https://img.shields.io/badge/Go-1.20-blue)

## 🚀 Installing

```sh
curl -L -o cloudcents https://github.com/cloudcents-ai/cli/raw/main/cloudcents
chmod +x cloudcents
echo 'export PATH=$PATH:'$(pwd) >> ~/.bash_profile
source ~/.bash_profile
```

## 🏁 Getting Started

![CLI Demo](img/cli.gif)

### 🔑 Authenticate with your API key
```sh
cloudcents auth <your-api-key>
```

### ✅ View and complete your task checklist
```
cloudcents checklist
```

### 💲 Get pricing for AWS, GCP, and Azure
```
cloudcents prices
```

### 💬 Chat with the Cloud Cents API
```
cloudcents chat "What is the advantage of AWS over GCP?"                  
```
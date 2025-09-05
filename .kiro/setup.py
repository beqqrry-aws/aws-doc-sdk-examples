from setuptools import setup, find_packages

setup(
    name="kiro-rag-server",
    version="0.1.0",
    py_modules=["rag_server"],  # This tells setup.py to include rag_server.py
    install_requires=[
        "fastmcp>=0.1.0",
        "langchain>=0.0.267",
        "langchain-community>=0.0.1",
        "sentence-transformers>=2.2.2",
        "chromadb>=0.4.13",
        "unstructured>=0.10.0",
        "pdf2image>=1.16.3",
        "pytesseract>=0.3.10",
        "docx2txt>=0.8",
    ],
    entry_points={
        "console_scripts": [
            "kiro-rag-server=rag_server:main",
        ],
    },
    python_requires=">=3.8",
)
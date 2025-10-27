from setuptools import setup, find_packages

setup(
    name='dag-utils',
    version='0.1.0',
    packages=find_packages(),
    install_requires=[
        'requests',
        'beautifulsoup4'
    ],
    python_requires='>=3.8',
)

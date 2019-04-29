package com.javabycode.springboot;

import com.amazonaws.auth.DefaultAWSCredentialsProviderChain;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.Bucket;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.SynchronousQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

@RestController
public class HelloWorldController {

    private final AmazonS3 s3;
    private AtomicInteger google = new AtomicInteger(0);
    private AtomicInteger facebook = new AtomicInteger(0);
    private ExecutorService googleExecutor;
    private ExecutorService facebookExecutor;

    public HelloWorldController() {
        s3 = AmazonS3ClientBuilder
            .standard()
            .withCredentials(new DefaultAWSCredentialsProviderChain())
            .build();
        googleExecutor = buildExecutorService(new SynchronousQueue<>());
        facebookExecutor = buildExecutorService(new ArrayBlockingQueue<>(250, true));
    }

    @RequestMapping("/s3")
    public String s3() {
        List<Bucket> buckets = s3.listBuckets();
        StringBuilder stringBuilder = new StringBuilder();
        stringBuilder.append("Your Amazon S3 buckets are:").append("\n");
        for (Bucket b : buckets) {
            stringBuilder.append("* ").append(b.getName()).append("\n");
        }

        return stringBuilder.toString();
    }

    @RequestMapping("/google")
    public CompletableFuture<String> google() {
        System.out.println("Google: " + google.addAndGet(1));
        CompletableFuture<String> googleFuture = new CompletableFuture<>();
        googleExecutor.submit(() -> {
            try {
                Thread.sleep(1000);
                googleFuture.complete("google");
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        });

        return googleFuture;
    }

    private ExecutorService buildExecutorService(BlockingQueue<Runnable> queue) {
        return new ThreadPoolExecutor(5, 5, 0L, TimeUnit.MILLISECONDS, queue);
    }

    @RequestMapping("/facebook")
    public CompletableFuture<String> facebook() {
        System.out.println("Facebook: " + facebook.addAndGet(1));
        CompletableFuture<String> facebookFuture = new CompletableFuture<>();
        facebookExecutor.submit(() -> {
            try {
                facebookFuture.complete(getRequest("https://www.facebook.com"));
            } catch (IOException e) {
                e.printStackTrace();
            }
        });

        return facebookFuture;
    }

    private String getRequest(String url) throws IOException {
        URL myUrl = new URL(url);
        HttpURLConnection con = (HttpURLConnection) myUrl.openConnection();
        con.setRequestMethod("GET");
        con.setConnectTimeout(10 * 1000);
        System.out.println(url + ":" + con.getResponseCode());
        StringBuilder content;

        try (BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()))) {
            String line;
            content = new StringBuilder();

            while ((line = in.readLine()) != null) {
                content.append(line);
                content.append(System.lineSeparator());
            }
        }
        con.disconnect();
        return content.toString();
    }
}

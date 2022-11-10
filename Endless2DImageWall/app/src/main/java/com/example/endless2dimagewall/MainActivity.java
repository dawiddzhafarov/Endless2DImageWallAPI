package com.example.endless2dimagewall;


import androidx.appcompat.app.AppCompatActivity;
import androidx.constraintlayout.widget.ConstraintLayout;

import android.annotation.SuppressLint;
import android.app.ProgressDialog;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.AsyncTask;
import android.os.Build;
import android.os.Bundle;
import android.util.Base64;
import android.util.Log;
import android.view.MotionEvent;
import android.view.View;
import android.view.Window;
import android.widget.FrameLayout;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.RelativeLayout;

import com.google.firebase.crashlytics.buildtools.reloc.org.apache.commons.io.IOUtils;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.lang.invoke.ConstantCallSite;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;
import java.nio.charset.Charset;
import java.sql.SQLOutput;
import java.util.Iterator;
import java.util.Scanner;
import java.lang.Object;
import java.util.concurrent.CompletableFuture;

import org.json.*;

import reactor.core.publisher.Mono;


public class MainActivity extends AppCompatActivity{

    float dx=0,dy=0,x=0,y=0;
    URL url;

    @SuppressLint("ClickableViewAccessibility")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        ConstraintLayout layout = (ConstraintLayout) findViewById(R.id.main);

        //new JsonTask().execute("http://10.0.2.2:8080/image");
        try {
            url = new URL("http://10.0.2.2:8080/images?z=1&x=1&y=1");
//            url = new URL("https://google.com");
        } catch (MalformedURLException e) {
            e.printStackTrace();
            url = null;
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.N) {
            CompletableFuture.supplyAsync(() -> getJson2(url)).thenAccept(this::testShowImg);
        }
        System.out.println("IIIIIIIIIIIIIIIIIIIII");

        layout.setOnTouchListener((v, event) -> {

            switch(event.getAction()) {
                case MotionEvent.ACTION_DOWN: {
                    x = event.getRawX();
                    y = event.getRawY();
                    dx = x - v.getX();
                    dy = y - v.getY();
                }
                break;
                case MotionEvent.ACTION_MOVE: {
                    v.setX(event.getRawX() - dx);
                    v.setY(event.getRawY() - dy);
                }
                break;
                case MotionEvent.ACTION_UP: {
                    //your stuff
                }
            }
            return true;
        });
    }
    public static JSONObject getJson(URL url) {
        String json = null;
        try {
            json = IOUtils.toString(url, Charset.forName("UTF-8"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        try {
            return new JSONObject(json);
        } catch (JSONException e) {
            e.printStackTrace();
        }
        return null;
    }

    public JSONObject getJson2(URL url) {
        HttpURLConnection connection = null;
        BufferedReader reader = null;
        try {
            connection = (HttpURLConnection) url.openConnection();
//            connection.addRequestProperty("username","sample_user");
//            connection.addRequestProperty("password","sample_password");
            String credential = Base64.encodeToString( ("sample_user"+":"+"sample_password").getBytes("UTF-8"), Base64.DEFAULT);
            connection.addRequestProperty("Authorization", "Basic " + credential.substring(0, credential.length()-1));
            System.out.println(connection);
            connection.connect();
            System.out.println(connection.getResponseCode());

            InputStream stream = connection.getInputStream();

            reader = new BufferedReader(new InputStreamReader(stream));

            StringBuffer buffer = new StringBuffer();
            String line = "";

            while ((line = reader.readLine()) != null) {
                buffer.append(line+"\n");
                Log.d("Response: ", "> " + line);   //here u ll get whole response...... :-)

            }
            return new JSONObject(buffer.toString());

        } catch (MalformedURLException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (JSONException e) {
            e.printStackTrace();
        } finally {
            if (connection != null) {
                connection.disconnect();
            }
            try {
                if (reader != null) {
                    reader.close();
                }
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
        return null;
    }

    public void testShowImg(JSONObject json) {
        try {
            JSONArray jsonArray = json.getJSONArray("images"); // wchodze głebiej, mam tablice 0,1,2
            String buffer = jsonArray.get(1).toString(); // pobieram element 1, zamieniam na string
            System.out.println(buffer);
            JSONObject jsonNested = new JSONObject(buffer); // zamieniam string na json
            System.out.println(jsonNested.get("base64"));
            String buffer2 = jsonNested.get("base64").toString();// pobieram element base64 z listy
            buffer = buffer2.substring(buffer2.indexOf(",") + 1); // usuwam niepotrzebne rzeczy
            buffer.trim();
            String encodedImage = buffer;
            byte[] decodedString = Base64.decode(encodedImage, Base64.DEFAULT);
            Bitmap decodedByte = BitmapFactory.decodeByteArray(decodedString, 0, decodedString.length);
            ImageView img = (ImageView) findViewById(R.id.imageView);
            img.setImageBitmap(decodedByte);
        }catch (JSONException e){}

    }

    private class JsonTask extends AsyncTask<String, String, String> {

        protected String doInBackground(String... params) {

            HttpURLConnection connection = null;
            BufferedReader reader = null;
            try {
                URL url = new URL(params[0]);
                connection = (HttpURLConnection) url.openConnection();
                connection.connect();;

                InputStream stream = connection.getInputStream();

                reader = new BufferedReader(new InputStreamReader(stream));

                StringBuffer buffer = new StringBuffer();
                String line = "";

                while ((line = reader.readLine()) != null) {
                    buffer.append(line+"\n");
                    Log.d("Response: ", "> " + line);   //here u ll get whole response...... :-)

                }

                return buffer.toString();

            } catch (MalformedURLException e) {
                e.printStackTrace();
            } catch (IOException e) {
                e.printStackTrace();
            } finally {
                if (connection != null) {
                    connection.disconnect();
                }
                try {
                    if (reader != null) {
                        reader.close();
                    }
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            return null;
        }

        @Override
        protected void onPostExecute(String result) {
            super.onPostExecute(result);
            System.out.println("AAAAAAAAAAAAAAAA");
            System.out.println(result);
            try {
                JSONObject json = new JSONObject(result); //zamiana ze stringa na jsona
                JSONArray jsonArray = json.getJSONArray("images"); // wchodze głebiej, mam tablice 0,1,2
                String buffer = jsonArray.get(1).toString(); // pobieram element 1, zamieniam na string
                System.out.println(buffer);
                JSONObject jsonNested = new JSONObject(buffer); // zamieniam string na json
                System.out.println(jsonNested.get("base64"));
                String buffer2 = jsonNested.get("base64").toString();// pobieram element base64 z listy
                buffer = buffer2.substring(buffer2.indexOf(",") + 1); // usuwam niepotrzebne rzeczy
                buffer.trim();
                String encodedImage = buffer;
                byte[] decodedString = Base64.decode(encodedImage, Base64.DEFAULT);
                Bitmap decodedByte = BitmapFactory.decodeByteArray(decodedString, 0, decodedString.length);
                ImageView img = (ImageView) findViewById(R.id.imageView);
                img.setImageBitmap(decodedByte);

            } catch (JSONException e) {
                e.printStackTrace();
            }
        }
    }


}



